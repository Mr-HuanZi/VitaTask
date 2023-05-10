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
	"github.com/duke-git/lancet/v2/maputil"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/gin-gonic/gin"
	"github.com/golang-module/carbon/v2"
	"gorm.io/gorm"
	"strings"
	"time"
)

type TaskLogLogic struct {
	Orm *gorm.DB
	ctx *gin.Context
}

func NewTaskLogLogic(ctx *gin.Context) *TaskLogLogic {
	return &TaskLogLogic{
		Orm: db.Db, // 赋予ORM实例
		ctx: ctx,   // 传递上下文
	}
}

// Add 新增日志
func (receiver TaskLogLogic) Add(data types.TaskLogForm) (*model.TaskLog, error) {
	// 检查任务是否存在
	// 此处会检索软删除的记录
	err := receiver.Orm.Model(&model.Task{}).Unscoped().First(&model.Task{}, data.TaskId).Error
	if err != nil {
		// 有错误就代表不存在
		return nil, err
	}
	// 操作类型是否合法
	if !slice.Contain(maputil.Keys(constant.GetTaskLogOperatorMaps()), data.OperateType) {
		return nil, exception.NewException(response.TaskOperatorTypeIllegal)
	}
	// 操作人
	if data.Operator <= 0 {
		// 获取当前用户
		currUser, err := user.CurrUser(receiver.ctx)
		if err != nil {
			return nil, err
		}
		data.Operator = currUser.ID
	} else {
		// 检查操作人是否存在
		_, err = NewMemberLogic(receiver.ctx).GetOne(data.Operator)
		if err != nil {
			return nil, err
		}
	}
	// 操作时间
	if data.OperateTime <= 0 {
		data.OperateTime = carbon.Now().TimestampMilli()
	}
	logData := &model.TaskLog{
		TaskId:      data.TaskId,
		OperateType: data.OperateType,
		Operator:    data.Operator,
		OperateTime: data.OperateTime,
		Message:     data.Message,
	}

	err = receiver.Orm.Create(logData).Error
	return logData, err
}

// Delete 删除单个日志
func (receiver TaskLogLogic) Delete(id uint64) error {
	return receiver.Orm.Model(&model.TaskLog{}).Delete(&model.TaskLog{ID: id}).Error
}

// Clear 清空某个任务的日志
func (receiver TaskLogLogic) Clear(taskId uint) error {
	return receiver.Orm.Where("task_id = ?", taskId).Delete(&model.TaskLog{}).Error
}

func (receiver TaskLogLogic) List(query types.TaskLogQuery) *types.PagedResult[model.TaskLog] {
	var (
		logs  []model.TaskLog
		count int64
	)
	tx := receiver.QueryHandle(query)
	if tx == nil {
		// tx 为空一般情况是提供了项目ID但是该项目下无任何任务
		return extend.EmptyPagedResult[model.TaskLog]()
	}
	// 获取数量
	tx.Count(&count)

	err := tx.Scopes(db.Paginate(&query.Page, &query.PageSize)).
		Preload("Task", func(db *gorm.DB) *gorm.DB {
			// 包括已删除的
			return db.Unscoped()
		}).
		Preload("OperatorInfo").
		Order("create_time DESC").Order("operate_time DESC").
		Find(&logs).Error

	if err != nil {
		logs = make([]model.TaskLog, 0)
		_ = exception.ErrorHandle(err, response.DbQueryError, "任务日志列表查询失败: ")
	}

	return extend.PagedResult[model.TaskLog](logs, count, int64(query.Page))
}

// QueryHandle 查询处理
func (receiver TaskLogLogic) QueryHandle(query types.TaskLogQuery) *gorm.DB {
	tx := receiver.Orm.Model(&model.TaskLog{})
	var taskIdsLimit []uint

	// 项目ID
	if query.ProjectId > 0 {
		receiver.Orm.Model(&model.Task{}).
			Unscoped().
			Select([]string{"id"}).
			Where("project_id = ?", query.ProjectId).
			Find(&taskIdsLimit)
		if len(taskIdsLimit) > 0 {
			query.TaskIds = append(taskIdsLimit, query.TaskIds...)
		} else {
			// 如果提供了项目ID参数，但是项目内任务列表为空，则直接返回nil
			return nil
		}
	}
	// 任务ID列表
	if len(query.TaskIds) > 0 {
		tx = tx.Where("task_id IN ?", query.TaskIds)
	}
	// 操作类型
	if len(strings.TrimSpace(query.OperateType)) > 0 {
		tx = tx.Where("operate_type = ?", query.OperateType)
	}
	// 操作人
	if query.Operator > 0 {
		tx = tx.Where("operator = ?", query.Operator)
	}
	// 时间范围
	if len(query.CreateTime) >= 2 {
		createTimeRange, err := time_tool.ParseStartEndTimeToUnix(query.CreateTime, time.DateTime, "milli")
		if err == nil {
			tx = tx.Where(
				"create_time BETWEEN ? AND ?",
				createTimeRange[0],
				createTimeRange[1],
			)
		}
	}
	// 计划时间范围
	if len(query.OperateTime) >= 2 {
		planTimeRange, err := time_tool.ParseStartEndTimeToUnix(query.OperateTime, time.DateTime, "milli")
		if err == nil {
			tx = tx.Where(
				"operate_time BETWEEN ? AND ?",
				planTimeRange[0],
				planTimeRange[1],
			)
		}
	}
	return tx
}

// Operators 获取操作类型
func (receiver TaskLogLogic) Operators() []map[string]string {
	operators := constant.GetTaskLogOperatorMaps()
	s := make([]map[string]string, len(operators))

	i := 0
	for k, item := range operators {
		s[i] = map[string]string{"label": item, "value": k}
		i++
	}
	return s
}
