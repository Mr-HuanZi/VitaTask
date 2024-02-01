package data

import (
	"VitaTaskGo/internal/api/model/dto"
	"VitaTaskGo/internal/repo"
	"VitaTaskGo/pkg/db"
	"VitaTaskGo/pkg/time_tool"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strings"
	"time"
)

type TaskLogRepo struct {
	tx  *gorm.DB
	ctx *gin.Context
}

func (r *TaskLogRepo) Create(data *repo.TaskLog) error {
	return r.tx.Create(&data).Error
}

func (r *TaskLogRepo) Save(data *repo.TaskLog) error {
	return r.tx.Save(&data).Error
}

func (r *TaskLogRepo) Delete(id uint64) error {
	return r.tx.Delete(&repo.TaskLog{}, id).Error
}

func (r *TaskLogRepo) Get(id uint64) (*repo.TaskLog, error) {
	var d *repo.TaskLog
	err := r.tx.First(&d, id).Error
	return d, err
}

func (r *TaskLogRepo) UpdateField(id uint64, field string, value interface{}) error {
	return r.tx.Model(&repo.TaskLog{}).Where("id = ?", id).Update(field, value).Error
}

func (r *TaskLogRepo) DeleteByTask(taskId uint) error {
	return r.tx.Where("task_id = ?", taskId).Delete(&repo.TaskLog{}).Error
}

func (r *TaskLogRepo) PageListTaskLog(query dto.TaskLogQuery) ([]repo.TaskLog, int64, error) {
	var (
		logs         []repo.TaskLog = nil
		total        int64
		taskIdsLimit []uint
	)

	tx := r.tx.Model(&repo.TaskLog{})

	// 项目ID
	if query.ProjectId > 0 {
		// todo 如何处理这一段？总不能实例化一个TaskRepo吧，虽然也没问题
		// todo 但是这样子好像也没问题呀
		r.tx.Model(&repo.Task{}).
			Unscoped().
			Select([]string{"id"}).
			Where("project_id = ?", query.ProjectId).
			Find(&taskIdsLimit)

		if len(taskIdsLimit) > 0 {
			query.TaskIds = append(taskIdsLimit, query.TaskIds...)
		} else {
			// 如果提供了项目ID参数，但是项目内任务列表为空，则直接返回nil
			return logs, total, nil
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

	// 获取总数
	err := tx.Count(&total).Error
	if err != nil {
		return logs, total, err
	}

	err = tx.Scopes(db.Paginate(&query.Page, &query.PageSize)).
		Preload("Task", func(db *gorm.DB) *gorm.DB {
			// 包括已删除的
			return db.Unscoped()
		}).
		Preload("OperatorInfo").
		Order("create_time DESC").Order("operate_time DESC").
		Find(&logs).Error

	return logs, total, err
}

func NewTaskLogRepo(tx *gorm.DB, ctx *gin.Context) repo.TaskLogRepo {
	return &TaskLogRepo{
		tx:  tx,
		ctx: ctx,
	}
}
