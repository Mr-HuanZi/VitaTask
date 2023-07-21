package service

import (
	"VitaTaskGo/internal/biz"
	"VitaTaskGo/internal/data"
	"VitaTaskGo/internal/dto"
	"VitaTaskGo/internal/pkg/db"
	"VitaTaskGo/internal/pkg/exception"
	"VitaTaskGo/internal/pkg/response"
	"github.com/duke-git/lancet/v2/maputil"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TaskGroupService struct {
	Orm  *gorm.DB
	ctx  *gin.Context
	repo biz.TaskGroupRepo
}

func NewTaskGroupService(tx *gorm.DB, ctx *gin.Context) *TaskGroupService {
	return &TaskGroupService{
		Orm:  tx,  // 赋予ORM实例
		ctx:  ctx, // 传递上下文
		repo: data.NewTaskGroupRepo(tx, ctx),
	}
}

// Add 新增组
func (receiver TaskGroupService) Add(taskGroupFormDto dto.TaskGroupForm) (*biz.TaskGroup, error) {
	// 判断项目是否存在
	if !data.NewProjectRepo(receiver.Orm, receiver.ctx).Exist(taskGroupFormDto.ProjectId) {
		return nil, exception.NewException(response.ProjectNotExist)
	}

	// 生成数据
	taskGroup := &biz.TaskGroup{
		ProjectId: taskGroupFormDto.ProjectId,
		Name:      taskGroupFormDto.Name,
	}
	// 保存
	err := receiver.repo.Create(taskGroup)
	return taskGroup, err
}

// Update 编辑组
func (receiver TaskGroupService) Update(taskGroupFormDto dto.TaskGroupForm) (*biz.TaskGroup, error) {
	// 判断项目是否存在
	if !data.NewProjectRepo(receiver.Orm, receiver.ctx).Exist(taskGroupFormDto.ProjectId) {
		return nil, exception.NewException(response.ProjectNotExist)
	}

	// 获取任务组数据
	taskGroup, err := receiver.repo.Get(taskGroupFormDto.ID)
	if err != nil {
		return nil, db.FirstQueryErrorHandle(err, response.TaskGroupNotExist)
	}

	// 更新字段
	taskGroup.ProjectId = taskGroupFormDto.ProjectId
	taskGroup.Name = taskGroupFormDto.Name
	// 保存
	err = receiver.repo.Save(taskGroup)
	return taskGroup, err
}

func (receiver TaskGroupService) Delete(groupId uint) error {
	// 任务组是否存在
	if !receiver.repo.Exist(groupId) {
		return exception.NewException(response.TaskGroupNotExist)
	}

	return exception.ErrorHandle(receiver.repo.Delete(groupId), response.DbExecuteError)
}

func (receiver TaskGroupService) List(query dto.TaskGroupQuery) dto.PagedResult[biz.TaskGroup] {
	list, total, err := receiver.repo.PageListTaskLog(query)

	if err != nil {
		list = make([]biz.TaskGroup, 0)
		_ = exception.ErrorHandle(err, response.DbQueryError, "任务组列表查询失败: ")
	}

	return dto.PagedResult[biz.TaskGroup]{
		Items: list,
		Total: total,
		Page:  int64(query.Page),
	}
}

// Detail 获取详情（带关联）
func (receiver TaskGroupService) Detail(groupId uint) (*biz.TaskGroup, error) {
	taskGroup, err := receiver.repo.Detail(groupId)
	if err != nil {
		return nil, db.FirstQueryErrorHandle(err, response.TaskGroupNotExist)
	}

	return taskGroup, nil
}

// SimpleList 简单列表
func (receiver TaskGroupService) SimpleList(projectId uint) []dto.UniversalSimpleList[uint] {
	list, err := receiver.repo.SimpleList(projectId)

	if err != nil {
		_ = exception.ErrorHandle(err, response.DbQueryError, "查询任务组简单列表失败: ")
		return nil
	}

	simpleList := make([]dto.UniversalSimpleList[uint], len(list))
	for i, s := range list {
		simpleList[i] = dto.UniversalSimpleList[uint]{
			Label: s.Name,
			Value: s.ID,
		}
	}

	return simpleList
}

// SimpleGroupList 简单列表-带项目分组
func (receiver TaskGroupService) SimpleGroupList() []dto.UniversalSimpleGroupList[uint] {
	var projectMap = make(map[uint]*biz.Project)

	list, err := receiver.repo.SimpleList(0)
	if err != nil {
		_ = exception.ErrorHandle(err, response.DbQueryError, "查询任务组简单列表失败: ")
		return nil
	}

	simpleGroupMap := make(map[uint]dto.UniversalSimpleGroupList[uint])

	for _, s := range list {
		label := ""
		if project, ok := projectMap[s.ProjectId]; ok {
			// 项目已缓存
			label = project.Name
		} else {
			// 项目未缓存
			project, err := data.NewProjectRepo(receiver.Orm, receiver.ctx).GetProject(s.ProjectId)
			if err != nil {
				continue
			}

			label = project.Name
			// 项目数据压入缓存
			projectMap[project.ID] = project
		}
		// 获取该项目的Options
		if group, ok := simpleGroupMap[s.ProjectId]; ok {
			// 加入已有的分组
			if len(group.Options) <= 0 {
				group.Options = []dto.UniversalSimpleList[uint]{{
					Label: s.Name,
					Value: s.ID,
				}}
			} else {
				group.Options = append(group.Options, dto.UniversalSimpleList[uint]{
					Label: s.Name,
					Value: s.ID,
				})
			}
			// 记得赋值回去
			simpleGroupMap[s.ProjectId] = group
		} else {
			// 新的分组
			simpleGroupMap[s.ProjectId] = dto.UniversalSimpleGroupList[uint]{
				Label: label,
				Options: []dto.UniversalSimpleList[uint]{{
					Label: s.Name,
					Value: s.ID,
				}},
			}
		}
	}

	return maputil.Values(simpleGroupMap)
}
