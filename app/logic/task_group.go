package logic

import (
	"VitaTaskGo/app/exception"
	"VitaTaskGo/app/extend/time_tool"
	"VitaTaskGo/app/model"
	"VitaTaskGo/app/response"
	"VitaTaskGo/app/types"
	"VitaTaskGo/library/db"
	"errors"
	"github.com/duke-git/lancet/v2/maputil"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"time"
)

type TaskGroupLogic struct {
	Db  *gorm.DB
	ctx *gin.Context
}

func NewTaskGroupLogic(ctx *gin.Context) *TaskGroupLogic {
	return &TaskGroupLogic{
		Db:  db.Db, // 赋予ORM实例
		ctx: ctx,   // 传递上下文
	}
}

// Add 新增组
func (receiver TaskGroupLogic) Add(data types.TaskGroupForm) (*model.TaskGroup, error) {
	// 判断项目是否存在
	_, projectErr := NewProjectLogic(receiver.ctx).GetOneProject(data.ProjectId)
	if projectErr != nil {
		return nil, projectErr
	}
	// 生成数据
	taskGroup := &model.TaskGroup{
		ProjectId: data.ProjectId,
		Name:      data.Name,
	}

	// 保存
	err := receiver.Db.Create(taskGroup).Error
	return taskGroup, err
}

// Update 编辑组
func (receiver TaskGroupLogic) Update(data types.TaskGroupForm) (*model.TaskGroup, error) {
	// 判断项目是否存在
	_, projectErr := NewProjectLogic(receiver.ctx).GetOneProject(data.ProjectId)
	if projectErr != nil {
		return nil, projectErr
	}

	// 获取任务组数据
	taskGroup, err := receiver.GetOne(data.ID)
	if err != nil {
		return nil, err
	}

	// 更新字段
	taskGroup.ProjectId = data.ProjectId
	taskGroup.Name = data.Name
	// 保存
	err = receiver.Db.Save(taskGroup).Error
	return taskGroup, err
}

func (receiver TaskGroupLogic) Delete(groupId uint) error {
	// 获取任务组数据
	taskGroup, err := receiver.GetOne(groupId)
	if err != nil {
		return err
	}

	return receiver.Db.Delete(taskGroup).Error
}

func (receiver TaskGroupLogic) List(query types.TaskGroupQuery) types.PagedResult[model.TaskGroup] {
	var (
		list  []model.TaskGroup
		count int64
	)

	tx := receiver.QueryHandle(query)
	// 获取总数
	tx.Count(&count)
	// 查询列表
	err := tx.Scopes(db.Paginate(&query.Page, &query.PageSize)).
		Preload("Project").
		Order("create_time DESC").
		Find(&list).Error

	if err != nil {
		list = make([]model.TaskGroup, 0)
		_ = exception.ErrorHandle(err, response.DbQueryError, "任务组列表查询失败: ")
	}

	return types.PagedResult[model.TaskGroup]{
		Items: list,
		Total: count,
		Page:  int64(query.Page),
	}
}

// Detail 获取详情（带关联）
func (receiver TaskGroupLogic) Detail(groupId uint) (*model.TaskGroup, error) {
	var taskGroup model.TaskGroup
	err := receiver.Db.Preload("Project").
		Order("create_time DESC").
		First(&taskGroup, groupId).Error

	if err != nil {
		return nil, exception.ErrorHandle(err, response.TaskGroupNotExist, "获取单条任务组记录失败: ")
	}
	if taskGroup.ID <= 0 {
		return nil, exception.NewException(response.TaskGroupNotExist)
	}

	return &taskGroup, nil
}

// GetOne 获取单条记录（无关联）
func (receiver TaskGroupLogic) GetOne(groupId uint) (*model.TaskGroup, error) {
	var taskGroup model.TaskGroup
	err := receiver.Db.First(&taskGroup, groupId).Error
	// 检查 ErrRecordNotFound 错误
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, exception.NewException(response.TaskGroupNotExist)
	}
	return &taskGroup, err
}

// QueryHandle 查询处理
func (receiver TaskGroupLogic) QueryHandle(query types.TaskGroupQuery) *gorm.DB {
	tx := receiver.Db.Model(model.TaskGroup{})

	if query.ProjectId > 0 {
		tx = tx.Where("project_id = ?", query.ProjectId)
	}
	// 时间范围
	if len(query.CreateTime) >= 2 {
		createTimeRange, err := time_tool.ParseStartEndTimeToUnix(query.CreateTime, time.DateOnly, "milli")
		if err == nil {
			tx = tx.Where(
				"create_time BETWEEN ? AND ?",
				createTimeRange[0],
				createTimeRange[1],
			)
		}
	}
	if query.Name != "" {
		tx = tx.Where("name LIKE ?", "%"+query.Name+"%")
	}
	return tx
}

// SimpleList 简单列表
func (receiver TaskGroupLogic) SimpleList(projectId uint) []types.UniversalSimpleList[uint] {
	var list []struct {
		Name string
		Id   uint
	}

	tx := receiver.Db.Model(&model.TaskGroup{}).Select("id", "name").Where("project_id", projectId)
	if err := tx.Scan(&list).Error; err != nil {
		_ = exception.ErrorHandle(err, response.DbQueryError, "查询任务组简单列表失败: ")
	}

	simpleList := make([]types.UniversalSimpleList[uint], len(list))

	for i, s := range list {
		simpleList[i] = types.UniversalSimpleList[uint]{
			Label: s.Name,
			Value: s.Id,
		}
	}

	return simpleList
}

// SimpleGroupList 简单列表-带项目分组
func (receiver TaskGroupLogic) SimpleGroupList() []types.UniversalSimpleGroupList[uint] {
	var list []struct {
		Name      string
		Id        uint
		ProjectId uint
	}
	var projectMap = make(map[uint]*model.Project)

	tx := receiver.Db.Model(&model.TaskGroup{}).Select("id", "name", "project_id")
	if err := tx.Scan(&list).Error; err != nil {
		_ = exception.ErrorHandle(err, response.DbQueryError, "查询任务组简单列表失败: ")
	}

	simpleGroupMap := make(map[uint]types.UniversalSimpleGroupList[uint])

	for _, s := range list {
		label := ""
		if project, ok := projectMap[s.ProjectId]; ok {
			// 项目已缓存
			label = project.Name
		} else {
			// 项目未缓存
			project, err := NewProjectLogic(receiver.ctx).GetOneProject(s.ProjectId)
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
				group.Options = []types.UniversalSimpleList[uint]{{
					Label: s.Name,
					Value: s.Id,
				}}
			} else {
				group.Options = append(group.Options, types.UniversalSimpleList[uint]{
					Label: s.Name,
					Value: s.Id,
				})
			}
			// 记得赋值回去
			simpleGroupMap[s.ProjectId] = group
		} else {
			// 新的分组
			simpleGroupMap[s.ProjectId] = types.UniversalSimpleGroupList[uint]{
				Label: label,
				Options: []types.UniversalSimpleList[uint]{{
					Label: s.Name,
					Value: s.Id,
				}},
			}
		}
	}

	return maputil.Values(simpleGroupMap)
}
