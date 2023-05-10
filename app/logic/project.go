package logic

import (
	"VitaTaskGo/app/constant"
	"VitaTaskGo/app/exception"
	"VitaTaskGo/app/extend"
	"VitaTaskGo/app/extend/time_tool"
	"VitaTaskGo/app/extend/user"
	"VitaTaskGo/app/model"
	"VitaTaskGo/app/response"
	"VitaTaskGo/app/types"
	"VitaTaskGo/library/db"
	"VitaTaskGo/library/state"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"time"
)

// ProjectLogic 项目模块逻辑
type ProjectLogic struct {
	Db  *gorm.DB // 事务DB实例
	ctx *gin.Context
}

func NewProjectLogic(ctx *gin.Context) *ProjectLogic {
	return &ProjectLogic{
		ctx: ctx,   // 传递上下文
		Db:  db.Db, // 赋予ORM实例
	}
}

// CreateProject 创建项目
func (receiver *ProjectLogic) CreateProject(name string, leaderUid uint64) (*model.Project, error) {
	// 判断负责人是否存在
	if !NewUserLogic(receiver.ctx).UserExist(leaderUid) {
		return nil, exception.NewException(response.ProjectLeaderNotExist)
	}
	// 创建新模型
	newProject := receiver.newProjectModel(name)
	// 获取当前用户
	currUser, err := user.CurrUser(receiver.ctx)
	if err != nil {
		return nil, err
	}

	// 在事务中进行
	transactionErr := receiver.Db.Transaction(func(tx *gorm.DB) error {
		// 创建项目
		if err := tx.Create(newProject).Error; err != nil {
			return err
		}

		projectMemberLogic := NewProjectMemberLogic(receiver.ctx)
		// 关联创建人
		if err := projectMemberLogic.Bind(newProject.ID, []uint64{currUser.ID}, constant.ProjectCreate); err != nil {
			return err
		}
		// 关联负责人
		return projectMemberLogic.Bind(newProject.ID, []uint64{leaderUid}, constant.ProjectLeader)
	})
	if err := exception.ErrorHandle(transactionErr, response.ProjectCreateFail, "创建项目失败: "); err != nil {
		return nil, err
	}

	// 重新读取
	return receiver.GetOneProject(newProject.ID)
}

// EditProject 更新项目
func (receiver *ProjectLogic) EditProject(projectId uint, name string, leader uint64) (*model.Project, error) {
	// 检查项目是否存在
	project, err := receiver.GetOneProject(projectId)
	if err != nil {
		return nil, err
	}
	// 在事务中进行
	transactionErr := receiver.Db.Transaction(func(tx *gorm.DB) error {
		// 改项目名
		project.Name = name
		// 保存项目名称
		if err := tx.Save(&project).Error; err != nil {
			return err
		}
		// 修改负责人
		if leader > 0 {
			// 判断负责人是否存在
			if !NewUserLogic(receiver.ctx).UserExist(leader) {
				return exception.NewException(response.ProjectLeaderNotExist)
			}
			// 获取该项目所有负责人
			projectMemberLogic := NewProjectMemberLogic(receiver.ctx)
			members, err := projectMemberLogic.GetMembersByRole(projectId, constant.ProjectLeader)
			if err != nil {
				return err
			}
			// 获取第一个
			member := members[0]
			// 赋予当前事务ORM实例
			projectMemberLogic.Db = tx
			// 开始移交
			return projectMemberLogic.Transfer(projectId, member.UserId, leader)
		}
		return nil
	})
	return project, exception.ErrorHandle(transactionErr, response.ProjectUpdateFail, "更新项目失败: ")
}

// GetProjectList 获取项目列表
func (receiver *ProjectLogic) GetProjectList(data types.ProjectListQuery, deleted bool) (*types.PagedResult[model.Project], error) {
	var (
		projects []model.Project
		count    int64
	)

	// 条件查询
	tx, err := receiver.QueryHandle(data, deleted)
	if err != nil {
		return nil, err
	}
	if deleted {
		// 回收站功能，只查询已删除的
		tx.Where("deleted_at IS NOT NULL")
	}

	// 计算总记录数
	tx.Model(&model.Project{}).Count(&count)
	// 获取记录
	if err := tx.Model(&model.Project{}).
		Scopes(db.Paginate(&data.Page, &data.PageSize)).
		Preload("Member.UserInfo").Find(&projects).Error; err != nil {
		projects = make([]model.Project, 0)
		_ = exception.ErrorHandle(err, response.ProjectNotExist, "项目列表查询失败: ")
		// 如果下面没有其它代码，可以继续往下执行
	}
	// 寻找负责人
	for i, project := range projects {
		for _, member := range project.Member {
			stateModifier := state.NewModifier(int(member.Role))
			if stateModifier.Exist(constant.ProjectLeader) {
				projects[i].Leader = member
				break
			}
		}
	}
	return &types.PagedResult[model.Project]{
		Items: projects,
		Total: count,
		Page:  int64(data.Page),
	}, nil
}

// GetSimpleList 获取简单的项目列表
func (receiver *ProjectLogic) GetSimpleList() []types.ProjectSimpleList {
	var simpleProjectList []types.ProjectSimpleList
	// 只检索当前用户所属的项目列表
	projectIds, err := receiver.MyProjectIds()
	if err != nil {
		_ = exception.ErrorHandle(err, response.DbQueryError)
		return nil
	}

	receiver.Db.Model(&model.Project{}).Where("id IN ?", projectIds).Find(&simpleProjectList)
	return simpleProjectList
}

// GetOneProject 获取单条项目记录
func (receiver *ProjectLogic) GetOneProject(projectId uint) (*model.Project, error) {
	var project *model.Project
	err := receiver.Db.Preload("Member.UserInfo").First(&project, projectId).Error
	if err != nil {
		return nil, exception.ErrorHandle(err, response.ProjectNotExist, "获取单条项目记录失败: ")
	}
	if project.ID <= 0 {
		return nil, exception.NewException(response.ProjectNotExist)
	}
	// 取得负责人
	for _, member := range project.Member {
		if member.Role == constant.ProjectLeader {
			project.Leader = member
			break
		}
	}
	return project, nil
}

// GetOneProjectByStringId 从字符串id中获取项目
func (receiver *ProjectLogic) GetOneProjectByStringId(id string) (*model.Project, error) {
	projectId := extend.ParseStringToUi64(id)

	// 查询项目
	return receiver.GetOneProject(uint(projectId))
}

// QueryHandle 查询条件处理
func (receiver *ProjectLogic) QueryHandle(data types.ProjectListQuery, deleted bool) (*gorm.DB, error) {
	var tx *gorm.DB
	if deleted {
		tx = receiver.Db.Unscoped() // 查询已删除的记录
	} else {
		tx = receiver.Db
	}

	// 只检索当前用户所属的项目列表
	projectIds, err := receiver.MyProjectIds()
	if err != nil {
		return nil, exception.ErrorHandle(err, response.DbQueryError)
	}
	tx = tx.Where("id IN ?", projectIds)

	if len(data.Time) >= 2 {
		timeRange, timeRangeErr := time_tool.ParseTimeRangeToUnix(data.Time, time.DateTime, "milli")
		if timeRangeErr != nil {
			return nil, exception.ErrorHandle(timeRangeErr, response.TimeParseFail)
		}
		tx = tx.Where("create_time BETWEEN ? AND ?", timeRange[0], timeRange[1])
	}
	if data.Name != "" {
		tx = tx.Where("name LIKE ?", "%"+data.Name+"%")
	}
	return tx, nil
}

// ProjectDelete 删除项目
func (receiver *ProjectLogic) ProjectDelete(projectId uint) error {
	project, err := receiver.GetOneProject(projectId)
	if err != nil {
		return err
	}
	// 执行删除
	tx := receiver.Db.Delete(&project)
	return exception.ErrorHandle(tx.Error, response.ProjectDeleteFail, "删除项目失败: ")
}

// ProjectArchive 项目归档
func (receiver *ProjectLogic) ProjectArchive(projectId uint) error {
	project, err := receiver.GetOneProject(projectId)
	if err != nil {
		return err
	}
	if project.Archive == constant.ProjectArchived {
		return exception.NewException(response.ProjectArchived) // 已归档
	}
	return exception.ErrorHandle(receiver.Db.Model(&project).Update("archive", constant.ProjectArchived).Error, response.ProjectArchiveFail, "项目归档失败: ")
}

// UnArchive 取消归档
func (receiver *ProjectLogic) UnArchive(projectId uint) error {
	project, err := receiver.GetOneProject(projectId)
	if err != nil {
		return err
	}
	if project.Archive == constant.ProjectNotArchive {
		return exception.NewException(response.ProjectNotArchived) // 未归档
	}
	return exception.ErrorHandle(receiver.Db.Model(&project).Update("archive", constant.ProjectNotArchive).Error, response.ProjectUnArchiveFail, "项目取消归档失败: ")
}

func (receiver *ProjectLogic) newProjectModel(name string) *model.Project {
	return &model.Project{
		Name:     name,
		Complete: 0,
		Archive:  0,
	}
}

// ProjectExist 项目是否存在
func (receiver *ProjectLogic) ProjectExist(projectId uint) bool {
	var count int64
	receiver.Db.Model(&model.ProjectMember{}).Where("project_id = ?", projectId).Count(&count)
	return count > 0
}

// Transfer 移交项目
func (receiver *ProjectLogic) Transfer(projectId uint, transferor, recipient uint64) error {
	// 如果两个id相同，不执行操作
	if transferor == recipient {
		return nil
	}
	// 判断项目是否存在
	if !receiver.ProjectExist(projectId) {
		return exception.NewException(response.ProjectNotExist)
	}
	// 检查是否在项目内
	projectMemberLogin := NewProjectMemberLogic(receiver.ctx)
	if !projectMemberLogin.InProject(projectId, transferor, constant.ProjectLeader) {
		return exception.NewException(response.MemberNotInProject)
	}
	// 开始移交
	return projectMemberLogin.Transfer(projectId, transferor, recipient)
}

// Archived 项目是否归档
func (receiver *ProjectLogic) Archived(projectId uint) bool {
	var count int64
	receiver.Db.Model(&model.Project{}).
		Where("id = ?", projectId).
		Where("archive = ?", 1).
		Count(&count)

	return count > 0
}

// MyProjects 获取当前用户所在项目的实例列表
func (receiver *ProjectLogic) MyProjects() ([]model.Project, error) {
	var (
		projectIds []uint
		projects   []model.Project
		err        error
	)

	projectIds, err = receiver.MyProjectIds()
	if err != nil {
		return nil, err
	}

	err = receiver.Db.Model(&model.Project{}).Where("id IN ?", projectIds).Find(&projects).Error
	return projects, err
}

// MyProjectIds 获取当前用户所在项目的ID列表
func (receiver *ProjectLogic) MyProjectIds() ([]uint, error) {
	var (
		projectIds []uint
		err        error
	)

	// 获取当前用户
	currUser, err := user.CurrUser(receiver.ctx)
	if err != nil {
		return nil, err
	}

	err = receiver.Db.Model(&model.ProjectMember{}).
		Select("project_id").
		Where("user_id = ?", currUser.ID).
		Group("project_id").
		Find(&projectIds).Error

	return projectIds, err
}
