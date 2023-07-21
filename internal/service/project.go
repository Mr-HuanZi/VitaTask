package service

import (
	"VitaTaskGo/internal/biz"
	"VitaTaskGo/internal/constant"
	"VitaTaskGo/internal/data"
	"VitaTaskGo/internal/dto"
	"VitaTaskGo/internal/pkg"
	"VitaTaskGo/internal/pkg/auth"
	"VitaTaskGo/internal/pkg/db"
	"VitaTaskGo/internal/pkg/exception"
	"VitaTaskGo/internal/pkg/response"
	"VitaTaskGo/internal/pkg/state"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// ProjectService 项目模块逻辑
type ProjectService struct {
	Db   *gorm.DB // 事务DB实例
	ctx  *gin.Context
	repo biz.ProjectRepo
}

func NewProjectService(tx *gorm.DB, ctx *gin.Context) *ProjectService {
	return &ProjectService{
		ctx:  ctx,   // 传递上下文
		Db:   db.Db, // 赋予ORM实例
		repo: data.NewProjectRepo(tx, ctx),
	}
}

// CreateProject 创建项目
func (receiver *ProjectService) CreateProject(name string, leaderUid uint64) (*biz.Project, error) {
	// 判断负责人是否存在
	if !data.NewUserRepo(receiver.Db, receiver.ctx).Exist(leaderUid) {
		return nil, exception.NewException(response.ProjectLeaderNotExist)
	}
	// 创建新模型
	newProject := receiver.newProjectModel(name)
	// 获取当前用户
	currUser, err := auth.CurrUser(receiver.ctx)
	if err != nil {
		return nil, err
	}

	// 在事务中进行
	transactionErr := receiver.Db.Transaction(func(tx *gorm.DB) error {
		// 创建项目
		if err := receiver.repo.CreateProject(newProject); err != nil {
			return err
		}

		projectMemberService := NewProjectMemberService(tx, receiver.ctx)
		// 关联创建人
		if err := projectMemberService.Bind(newProject.ID, []uint64{currUser.ID}, constant.ProjectCreate); err != nil {
			return err
		}
		// 关联负责人
		return projectMemberService.Bind(newProject.ID, []uint64{leaderUid}, constant.ProjectLeader)
	})
	if err := exception.ErrorHandle(transactionErr, response.ProjectCreateFail, "创建项目失败: "); err != nil {
		return nil, err
	}

	// 重新读取
	return receiver.GetOneProject(newProject.ID)
}

// EditProject 更新项目
func (receiver *ProjectService) EditProject(projectId uint, name string, leader uint64) (*biz.Project, error) {
	// 检查项目是否存在
	project, err := receiver.GetOneProject(projectId)
	if err != nil {
		return nil, err
	}
	// 在事务中进行
	transactionErr := receiver.Db.Transaction(func(tx *gorm.DB) error {
		// 创建新的Repo
		projectRepo := data.NewProjectRepo(tx, receiver.ctx)
		userRepo := data.NewUserRepo(tx, receiver.ctx)
		// 改项目名
		project.Name = name
		// 保存项目名称
		if err := projectRepo.SaveProject(project); err != nil {
			return err
		}

		// 修改负责人
		if leader > 0 {
			// 判断负责人是否存在
			if !userRepo.Exist(leader) {
				return exception.NewException(response.ProjectLeaderNotExist)
			}

			projectMemberService := NewProjectMemberService(tx, receiver.ctx)
			// 移交项目给新负责人， 如果新旧负责人ID相同则不执行任何操作
			return projectMemberService.Transfer(projectId, leader)
		}
		return nil
	})
	return project, exception.ErrorHandle(transactionErr, response.ProjectUpdateFail, "更新项目失败: ")
}

// GetProjectList 获取项目列表
func (receiver *ProjectService) GetProjectList(query dto.ProjectListQuery) (*dto.PagedResult[biz.Project], error) {
	// 只检索当前用户所属的项目列表
	projectIds, err := receiver.MyProjectIds()
	if err != nil {
		return nil, exception.ErrorHandle(err, response.DbQueryError)
	}

	projects, total, err := receiver.repo.PageListProject(query, projectIds)
	if err != nil {
		return pkg.PagedResult(projects, total, int64(query.Page)), err
	}

	// 寻找负责人
	for i, project := range projects {
		// todo 封装一下
		for _, member := range project.Member {
			stateModifier := state.NewModifier(int(member.Role))
			if stateModifier.Exist(constant.ProjectLeader) {
				projects[i].Leader = member
				break
			}
		}
	}
	// 返回结果
	return pkg.PagedResult(projects, total, int64(query.Page)), nil
}

// GetSimpleList 获取简单的项目列表
func (receiver *ProjectService) GetSimpleList() []dto.ProjectSimpleList {
	// 只检索当前用户所属的项目列表
	projectIds, err := receiver.MyProjectIds()
	if err != nil {
		logrus.Errorln("简单的项目列表获取当前用户所属项目失败：", err)
		return nil
	}

	simpleProjectList, err := receiver.repo.SimpleList(projectIds)
	if err != nil {
		logrus.Errorln("简单的项目列表数据获取失败：", err)
		return nil
	}

	return simpleProjectList
}

// GetOneProject 获取单条项目记录
func (receiver *ProjectService) GetOneProject(projectId uint) (*biz.Project, error) {
	project, err := receiver.repo.PreloadGetProject(projectId)
	if err != nil {
		return nil, db.FirstQueryErrorHandle(err, response.ProjectNotExist)
	}

	// 取得负责人
	// todo 封装一下
	for _, member := range project.Member {
		stateModifier := state.NewModifier(int(member.Role))
		if stateModifier.Exist(constant.ProjectLeader) {
			project.Leader = member
			break
		}
	}
	return project, nil
}

// GetOneProjectByStringId 从字符串id中获取项目
func (receiver *ProjectService) GetOneProjectByStringId(id string) (*biz.Project, error) {
	projectId := pkg.ParseStringToUi64(id)

	// 查询项目
	return receiver.GetOneProject(uint(projectId))
}

// ProjectDelete 删除项目
func (receiver *ProjectService) ProjectDelete(projectId uint) error {
	_, err := receiver.GetOneProject(projectId)
	if err != nil {
		return err
	}
	// 执行删除
	err = receiver.repo.DeleteProject(projectId)
	return exception.ErrorHandle(err, response.ProjectDeleteFail, "删除项目失败: ")
}

// ProjectArchive 项目归档
func (receiver *ProjectService) ProjectArchive(projectId uint) error {
	project, err := receiver.GetOneProject(projectId)
	if err != nil {
		return err
	}

	if project.Archive == constant.ProjectArchived {
		return exception.NewException(response.ProjectArchived) // 已归档
	}

	return exception.ErrorHandle(receiver.repo.UpdateField(project.ID, "archive", constant.ProjectArchived), response.ProjectArchiveFail, "项目归档失败: ")
}

// UnArchive 取消归档
func (receiver *ProjectService) UnArchive(projectId uint) error {
	project, err := receiver.GetOneProject(projectId)
	if err != nil {
		return err
	}

	if project.Archive == constant.ProjectNotArchive {
		return exception.NewException(response.ProjectNotArchived) // 未归档
	}

	return exception.ErrorHandle(receiver.repo.UpdateField(project.ID, "archive", constant.ProjectNotArchive), response.ProjectUnArchiveFail, "项目取消归档失败: ")
}

func (receiver *ProjectService) newProjectModel(name string) *biz.Project {
	return &biz.Project{
		Name:     name,
		Complete: 0,
		Archive:  0,
	}
}

// Transfer 移交项目
func (receiver *ProjectService) Transfer(projectId uint, recipient uint64) error {
	// 判断项目是否存在
	if !receiver.repo.Exist(projectId) {
		return exception.NewException(response.ProjectNotExist)
	}

	// 移交项目给新负责人， 如果新旧负责人ID相同则不执行任何操作
	return NewProjectMemberService(receiver.Db, receiver.ctx).Transfer(projectId, recipient)
}

// MyProjects 获取当前用户所在项目的实例列表
func (receiver *ProjectService) MyProjects() ([]biz.Project, error) {
	// 获取当前用户
	currUser, err := auth.CurrUser(receiver.ctx)
	if err != nil {
		return nil, err
	}

	projects, err := receiver.repo.GetUserProjects(currUser.ID)
	return projects, exception.ErrorHandle(err, response.DbQueryError)
}

// MyProjectIds 获取当前用户所在项目的ID列表
func (receiver *ProjectService) MyProjectIds() ([]uint, error) {
	projects, err := receiver.MyProjects()
	if err != nil {
		return nil, err
	}

	projectIds := make([]uint, len(projects))
	for i, project := range projects {
		projectIds[i] = project.ID
	}

	return projectIds, nil
}
