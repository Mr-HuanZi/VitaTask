package data

import (
	"VitaTaskGo/internal/biz"
	"VitaTaskGo/internal/constant"
	"VitaTaskGo/internal/dto"
	"VitaTaskGo/internal/pkg/db"
	"VitaTaskGo/internal/pkg/exception"
	"VitaTaskGo/internal/pkg/response"
	"VitaTaskGo/internal/pkg/time_tool"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"time"
)

type TaskRepo struct {
	tx  *gorm.DB
	ctx *gin.Context
}

func (r *TaskRepo) Create(data *biz.Task) error {
	return r.tx.Create(&data).Error
}

func (r *TaskRepo) Save(data *biz.Task) error {
	return r.tx.Save(&data).Error
}

func (r *TaskRepo) Delete(id uint) error {
	return r.tx.Delete(&biz.Task{}, id).Error
}

func (r *TaskRepo) Get(id uint) (*biz.Task, error) {
	var d *biz.Task
	err := r.tx.First(&d, id).Error
	return d, err
}

func (r *TaskRepo) UpdateField(id uint, field string, value interface{}) error {
	return r.tx.Model(&biz.Task{}).Where("id = ?", id).Update(field, value).Error
}

func (r *TaskRepo) UpdateFields(id uint, values interface{}) error {
	return r.tx.Model(&biz.Task{}).Where("id = ?", id).Updates(values).Error
}

func (r *TaskRepo) PageListProject(query dto.TaskListQueryBO) ([]biz.Task, int64, error) {
	var (
		list  []biz.Task = nil
		total int64
	)

	tx := r.tx.Model(biz.Task{})

	// 查询已删除的记录
	// 回收站功能，只查询已删除的
	if query.Deleted {
		tx = tx.Unscoped().Where("deleted_at IS NOT NULL")
	}

	if len(query.ProjectIds) > 0 {
		tx = tx.Where("project_id IN ?", query.ProjectIds)
	}

	// 任务组搜索
	if query.GroupId > 0 {
		tx = tx.Where("group_id = ?", query.GroupId)
	}

	// 负责人
	if len(query.LeaderTaskIds) > 0 {
		tx = tx.Where("id IN ?", query.LeaderTaskIds)
	}
	// 协助人
	if len(query.CollaboratorTaskIds) > 0 {
		tx = tx.Where("id IN ?", query.CollaboratorTaskIds)
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
	// 计划时间范围
	if len(query.PlanTime) >= 2 {
		planTimeRange, err := time_tool.ParseStartEndTimeToUnix(query.PlanTime, time.DateOnly, "milli")
		if err == nil {
			// todo 改为 "(start_date BETWEEN ? AND ?) OR (end_date BETWEEN ? AND ?)"是否会更合理一些？
			tx = tx.Where(
				"start_date >= ? AND end_date <= ?",
				planTimeRange[0],
				planTimeRange[1],
			)
		}
	}
	// 标题搜索
	titleQuery := ""
	if query.Title != "" {
		titleQuery = query.Title
	} else if query.Name != "" {
		titleQuery = query.Name
	}
	if titleQuery != "" {
		tx = tx.Where("title LIKE ?", "%"+titleQuery+"%")
	}

	// 计算总记录数
	err := tx.Model(&biz.Task{}).Count(&total).Error
	if err != nil {
		return list, 0, err
	}

	// 查询记录
	err = tx.Model(&biz.Task{}).
		Scopes(db.Paginate(&query.Page, &query.PageSize)).
		Preload("Project").
		Preload("Group").
		Preload("Member.UserInfo").
		Order("status ASC").Order("level DESC").Order("create_time DESC").
		Find(&list).Error

	return list, total, exception.ErrorHandle(err, response.DbQueryError)
}

func (r *TaskRepo) Detail(id uint) (*biz.Task, error) {
	var task *biz.Task
	err := r.tx.Preload("Project").
		Preload("Group").
		Preload("Member.UserInfo").
		Order("status ASC").Order("level DESC").Order("create_time DESC").
		First(&task, id).Error

	return task, err
}

func (r *TaskRepo) TaskNumber(projectId uint, status []int) (int64, error) {
	var count int64

	tx := r.tx.Model(&biz.Task{}).Where(&biz.Task{ProjectId: projectId})
	if len(status) > 0 {
		tx = tx.Where("status IN ?", status)
	}

	err := tx.Count(&count).Error

	return count, err
}

func (r *TaskRepo) GetTasksByProject(projectId uint, status []int) ([]biz.Task, error) {
	var list []biz.Task

	tx := r.tx.Model(&biz.Task{}).Where(&biz.Task{ProjectId: projectId})
	if len(status) > 0 {
		tx = tx.Where("status IN ?", status)
	}

	err := tx.Find(&list).Error

	return list, err
}

func (r *TaskRepo) CompletedQuantity(projectId uint, completeTime []int64) (int64, error) {
	var count int64

	tx := r.tx.Model(&biz.Task{}).Where(&biz.Task{ProjectId: projectId}).Where("status IN ?", []int{constant.TaskStatusCompleted, constant.TaskStatusArchived})
	if len(completeTime) >= 2 {
		tx = tx.Where("complete_date BETWEEN ? AND ?", completeTime[0], completeTime[1])
	}

	err := tx.Count(&count).Error

	return count, err
}

func (r *TaskRepo) CreatedQuantity(projectId uint, createTime []int64) (int64, error) {
	var count int64

	tx := r.tx.Model(&biz.Task{}).Where(&biz.Task{ProjectId: projectId})
	if len(createTime) >= 2 {
		tx = tx.Where("complete_date BETWEEN ? AND ?", createTime[0], createTime[1])
	}

	err := tx.Count(&count).Error

	return count, err
}

func NewTaskRepo(tx *gorm.DB, ctx *gin.Context) biz.TaskRepo {
	return &TaskRepo{
		tx:  tx,
		ctx: ctx,
	}
}
