package data

import (
	"VitaTaskGo/internal/biz"
	"VitaTaskGo/internal/dto"
	"VitaTaskGo/internal/pkg/db"
	"VitaTaskGo/internal/pkg/time_tool"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"time"
)

type TaskGroupRepo struct {
	tx  *gorm.DB
	ctx *gin.Context
}

func (r *TaskGroupRepo) Create(data *biz.TaskGroup) error {
	return r.tx.Create(&data).Error
}

func (r *TaskGroupRepo) Save(data *biz.TaskGroup) error {
	return r.tx.Save(&data).Error
}

func (r *TaskGroupRepo) Delete(id uint) error {
	return r.tx.Delete(&biz.TaskGroup{}, id).Error
}

func (r *TaskGroupRepo) Get(id uint) (*biz.TaskGroup, error) {
	var d *biz.TaskGroup
	err := r.tx.First(&d, id).Error
	return d, err
}

func (r *TaskGroupRepo) UpdateField(id uint, field string, value interface{}) error {
	return r.tx.Model(&biz.TaskGroup{}).Where("id = ?", id).Update(field, value).Error
}

func (r *TaskGroupRepo) Exist(id uint) bool {
	return r.tx.Select("id").Where("id = ?", id).First(&biz.TaskGroup{}).Error == nil
}

func (r *TaskGroupRepo) PageListTaskLog(query dto.TaskGroupQuery) ([]biz.TaskGroup, int64, error) {
	var (
		list  []biz.TaskGroup = nil
		total int64
	)

	tx := r.tx.Model(biz.TaskGroup{})

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

	// 获取总数
	err := tx.Count(&total).Error
	if err != nil {
		return list, total, err
	}

	// 查询列表
	err = tx.Scopes(db.Paginate(&query.Page, &query.PageSize)).
		Preload("Project").
		Order("create_time DESC").
		Find(&list).Error

	return list, total, err
}

func (r *TaskGroupRepo) Detail(id uint) (*biz.TaskGroup, error) {
	var d *biz.TaskGroup
	err := r.tx.Preload("Project").Order("create_time DESC").First(&d, id).Error
	return d, err
}

// SimpleList 简单列表
// 允许不提供项目ID
func (r *TaskGroupRepo) SimpleList(projectId uint) ([]biz.TaskGroup, error) {
	var l []biz.TaskGroup
	tx := r.tx.Model(&biz.TaskGroup{}).Select("id", "name", "project_id")
	// 如果不提供项目id
	if projectId <= 0 {
		tx = tx.Where("project_id", projectId)
	}
	err := tx.Find(&l).Error

	return l, err
}

func NewTaskGroupRepo(tx *gorm.DB, ctx *gin.Context) biz.TaskGroupRepo {
	return &TaskGroupRepo{
		tx:  tx,
		ctx: ctx,
	}
}
