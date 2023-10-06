package data

import (
	"VitaTaskGo/internal/biz"
	"VitaTaskGo/internal/dto"
	"VitaTaskGo/internal/pkg/db"
	"VitaTaskGo/internal/pkg/exception"
	"VitaTaskGo/internal/pkg/response"
	"VitaTaskGo/internal/pkg/time_tool"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"time"
)

type WorkflowRepo struct {
	tx  *gorm.DB
	ctx *gin.Context
}

func NewWorkflowRepo(tx *gorm.DB, ctx *gin.Context) biz.WorkflowRepo {
	return &WorkflowRepo{
		tx:  tx,
		ctx: ctx,
	}
}

func (r *WorkflowRepo) Create(data *biz.Workflow) error {
	return r.tx.Create(&data).Error
}

func (r *WorkflowRepo) Save(data *biz.Workflow) error {
	return r.tx.Save(&data).Error
}

func (r *WorkflowRepo) Delete(id uint) error {
	return r.tx.Delete(&biz.Workflow{}, id).Error
}

func (r *WorkflowRepo) Get(id uint) (*biz.Workflow, error) {
	var d *biz.Workflow
	err := r.tx.First(&d, id).Error
	return d, err
}

func (r *WorkflowRepo) UpdateField(id uint, field string, value interface{}) error {
	return r.tx.Model(&biz.Workflow{}).Where("id = ?", id).Update(field, value).Error
}

func (r *WorkflowRepo) UpdateFields(id uint, values interface{}) error {
	return r.tx.Model(&biz.Workflow{}).Where("id = ?", id).Updates(values).Error
}

func (r *WorkflowRepo) PageList(query dto.WorkflowListQueryDto) ([]biz.Workflow, int64, error) {
	var (
		list  []biz.Workflow = nil
		total int64
	)

	tx := r.tx.Model(biz.Workflow{})

	if query.ID > 0 {
		tx = tx.Where("id = ?", query.ID)
	}

	if len(query.TypeId) > 0 {
		tx = tx.Where("type_id IN ?", query.TypeId)
	}

	if len(query.CreateTime) >= 2 {
		// 将字符串时间转换为数字时间戳
		createTimeRange, err := time_tool.ParseStartEndTimeToUnix(query.CreateTime, time.DateTime, "milli")
		if err != nil {
			return nil, 0, exception.ErrorHandle(err, response.TimeParseFail)
		}
		tx = tx.Where(
			"create_time BETWEEN ? AND ?",
			createTimeRange[0],
			createTimeRange[1],
		)
	}

	// 工作流状态 的最小值是 0
	// 使用字符串类型是因为整型在初始化的时候为 0 值，和 工作流状态冲突
	if len(query.Status) > 0 {
		tx = tx.Where("status = ?", query.Status)
	}

	if len(query.Title) > 0 {
		tx = tx.Where("title LIKE ?", "%"+query.Title+"%")
	}

	if len(query.Serials) > 0 {
		tx = tx.Where("serials LIKE ?", "%"+query.Serials+"%")
	}

	// 查询已删除的记录
	if query.Deleted {
		tx = tx.Unscoped().Where("deleted_at IS NOT NULL")
	}

	// 计算总记录数
	err := tx.Count(&total).Error
	if err != nil {
		return list, 0, exception.ErrorHandle(err, response.DbQueryError)
	}

	// 查询记录
	err = tx.Scopes(db.Paginate(&query.Page, &query.PageSize)).
		// 关联当前操作人
		Preload("Operator").
		Order("create_time DESC").
		Find(&list).Error

	return list, total, exception.ErrorHandle(err, response.DbQueryError)
}

func (r *WorkflowRepo) GetDayTotal(start, end int64) (int64, error) {
	var count int64
	err := r.tx.Model(&biz.Workflow{}).Where("create_time BETWEEN ? AND ?", start, end).Count(&count).Error
	return count, err
}

func (r *WorkflowRepo) SetDbInstance(tx *gorm.DB) {
	r.tx = tx
}
