package data

import (
	"VitaTaskGo/internal/api/model/dto"
	"VitaTaskGo/internal/repo"
	"VitaTaskGo/pkg/db"
	"VitaTaskGo/pkg/exception"
	"VitaTaskGo/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type WorkflowLogRepo struct {
	tx  *gorm.DB
	ctx *gin.Context
}

func NewWorkflowLogRepo(tx *gorm.DB, ctx *gin.Context) repo.WorkflowLogRepo {
	return &WorkflowLogRepo{
		tx:  tx,
		ctx: ctx,
	}
}

func (r *WorkflowLogRepo) Create(data *repo.WorkflowLog) error {
	return r.tx.Create(&data).Error
}

func (r *WorkflowLogRepo) Save(data *repo.WorkflowLog) error {
	return r.tx.Save(&data).Error
}

func (r *WorkflowLogRepo) Delete(id uint) error {
	return r.tx.Delete(&repo.WorkflowLog{}, id).Error
}

func (r *WorkflowLogRepo) Get(id uint) (*repo.WorkflowLog, error) {
	var d *repo.WorkflowLog
	err := r.tx.First(&d, id).Error
	return d, err
}

func (r *WorkflowLogRepo) UpdateField(id uint, field string, value interface{}) error {
	return r.tx.Model(&repo.WorkflowLog{}).Where("id = ?", id).Update(field, value).Error
}

func (r *WorkflowLogRepo) UpdateFields(id uint, values interface{}) error {
	return r.tx.Model(&repo.WorkflowLog{}).Where("id = ?", id).Updates(values).Error
}

func (r *WorkflowLogRepo) PageList(query dto.WorkflowLogQueryBo) ([]repo.WorkflowLog, int64, error) {
	var (
		list  []repo.WorkflowLog = nil
		total int64
	)

	tx := r.tx.Model(repo.WorkflowLog{})

	if query.ID > 0 {
		tx = tx.Where("id = ?", query.ID)
	}

	if query.WorkflowId > 0 {
		tx = tx.Where("workflow_id = ?", query.WorkflowId)
	}

	if query.Node > 0 {
		tx = tx.Where("node = ?", query.Node)
	}

	if query.Operator > 0 {
		tx = tx.Where("operator = ?", query.Operator)
	}

	if len(query.CreateTime) >= 2 {
		tx = tx.Where(
			"create_time BETWEEN ? AND ?",
			query.CreateTime[0],
			query.CreateTime[1],
		)
	}

	if len(query.Action) > 0 {
		tx = tx.Where("action = ?", query.Action)
	}

	// 计算总记录数
	err := tx.Count(&total).Error
	if err != nil {
		return list, 0, exception.ErrorHandle(err, response.DbQueryError)
	}

	// 查询记录
	err = tx.Scopes(db.Paginate(&query.Page, &query.PageSize)).
		Order("create_time DESC").
		Find(&list).Error

	return list, total, exception.ErrorHandle(err, response.DbQueryError)
}

func (r *WorkflowLogRepo) SetDbInstance(tx *gorm.DB) {
	r.tx = tx
}
