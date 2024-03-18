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

type WorkflowNodeRepo struct {
	tx  *gorm.DB
	ctx *gin.Context
}

func NewWorkflowNodeRepo(tx *gorm.DB, ctx *gin.Context) repo.WorkflowNodeRepo {
	return &WorkflowNodeRepo{
		tx:  tx,
		ctx: ctx,
	}
}

func (r *WorkflowNodeRepo) Create(data *repo.WorkflowNode) error {
	return r.tx.Create(&data).Error
}

func (r *WorkflowNodeRepo) Save(data *repo.WorkflowNode) error {
	return r.tx.Save(&data).Error
}

func (r *WorkflowNodeRepo) Delete(id uint) error {
	return r.tx.Delete(&repo.WorkflowNode{}, id).Error
}

func (r *WorkflowNodeRepo) Get(id uint) (*repo.WorkflowNode, error) {
	var d *repo.WorkflowNode
	err := r.tx.First(&d, id).Error
	return d, err
}

func (r *WorkflowNodeRepo) UpdateField(id uint, field string, value interface{}) error {
	return r.tx.Model(&repo.WorkflowNode{}).Where("id = ?", id).Update(field, value).Error
}

func (r *WorkflowNodeRepo) UpdateFields(id uint, values interface{}) error {
	return r.tx.Model(&repo.WorkflowNode{}).Where("id = ?", id).Updates(values).Error
}

func (r *WorkflowNodeRepo) PageList(query dto.WorkflowNodeQueryBo) ([]repo.WorkflowNode, int64, error) {
	var (
		list  []repo.WorkflowNode = nil
		total int64
	)

	tx := r.tx.Model(repo.WorkflowNode{})

	if query.ID > 0 {
		tx = tx.Where("id = ?", query.ID)
	}

	if query.TypeId > 0 {
		tx = tx.Where("type_id = ?", query.TypeId)
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

	if len(query.Name) > 0 {
		tx = tx.Where("name LIKE ?", "%"+query.Name+"%")
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
		Order("create_time DESC").
		Find(&list).Error

	return list, total, exception.ErrorHandle(err, response.DbQueryError)
}

func (r *WorkflowNodeRepo) GetAppointNode(typeId uint, node int) (*repo.WorkflowNode, error) {
	var d *repo.WorkflowNode
	// struct查询会忽略0值，而工作流结束后的node字段就是0值
	err := r.tx.Model(&repo.WorkflowNode{}).Where(&repo.WorkflowNode{TypeId: typeId}).Where("node = ?", node).First(&d).Error
	return d, err
}

func (r *WorkflowNodeRepo) GetNextNode(typeId uint, currNode int) (*repo.WorkflowNode, error) {
	var d *repo.WorkflowNode
	err := r.tx.Model(&repo.WorkflowNode{}).
		Where(&repo.WorkflowNode{TypeId: typeId}).
		Where("node > ?", currNode).
		Order("node").
		First(&d).Error
	return d, err
}

func (r *WorkflowNodeRepo) FirstNode(typeId uint) (*repo.WorkflowNode, error) {
	var d *repo.WorkflowNode
	err := r.tx.Model(&repo.WorkflowNode{}).Where(&repo.WorkflowNode{TypeId: typeId}).Order("node ASC").First(&d).Error
	return d, err
}

func (r *WorkflowNodeRepo) GetTypeAll(typeId uint) ([]repo.WorkflowNode, error) {
	var list []repo.WorkflowNode

	err := r.tx.Model(&repo.WorkflowNode{}).
		Where(&repo.WorkflowNode{TypeId: typeId}).
		Order("node").
		Find(&list).Error

	return list, err
}

func (r *WorkflowNodeRepo) SetDbInstance(tx *gorm.DB) {
	r.tx = tx
}
