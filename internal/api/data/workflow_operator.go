package data

import (
	"VitaTaskGo/internal/repo"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type WorkflowOperatorRepo struct {
	tx  *gorm.DB
	ctx *gin.Context
}

func NewWorkflowOperatorRepo(tx *gorm.DB, ctx *gin.Context) repo.WorkflowOperatorRepo {
	return &WorkflowOperatorRepo{
		tx:  tx,
		ctx: ctx,
	}
}

func (r *WorkflowOperatorRepo) Create(data *repo.WorkflowOperator) error {
	return r.tx.Create(&data).Error
}

func (r *WorkflowOperatorRepo) Save(data *repo.WorkflowOperator) error {
	return r.tx.Save(&data).Error
}

func (r *WorkflowOperatorRepo) Delete(id uint) error {
	return r.tx.Delete(&repo.WorkflowOperator{}, id).Error
}

func (r *WorkflowOperatorRepo) Get(id uint) (*repo.WorkflowOperator, error) {
	var d *repo.WorkflowOperator
	err := r.tx.First(&d, id).Error
	return d, err
}

func (r *WorkflowOperatorRepo) UpdateField(id uint, field string, value interface{}) error {
	return r.tx.Model(&repo.WorkflowOperator{}).Where("id = ?", id).Update(field, value).Error
}

func (r *WorkflowOperatorRepo) UpdateFields(id uint, values interface{}) error {
	return r.tx.Model(&repo.WorkflowOperator{}).Where("id = ?", id).Updates(values).Error
}

func (r *WorkflowOperatorRepo) GetWorkflowOperatorByNode(workflowId uint, node int) ([]repo.WorkflowOperator, error) {
	var l []repo.WorkflowOperator
	err := r.tx.Model(&repo.WorkflowOperator{}).Where(&repo.WorkflowOperator{WorkflowId: workflowId, Node: node}).Find(&l).Error
	return l, err
}

func (r *WorkflowOperatorRepo) SetDbInstance(tx *gorm.DB) {
	r.tx = tx
}

func (r *WorkflowOperatorRepo) OtherOperator(workflowId uint, node int, userId uint64) (bool, error) {
	var count int64
	err := r.tx.Model(&repo.WorkflowOperator{}).
		Where(&repo.WorkflowOperator{WorkflowId: workflowId, Node: node}).
		Where("handled = ?", 0).
		Where("user_id <> ?", userId).
		Count(&count).Error

	return count > 0, err
}

func (r *WorkflowOperatorRepo) RemoveWorkflowAllOperator(workflowId uint) error {
	return r.tx.Where(&repo.WorkflowOperator{WorkflowId: workflowId}).Delete(&repo.WorkflowOperator{}).Error
}

func (r *WorkflowOperatorRepo) SetHandled(workflowId uint, node int, userId uint64) error {
	return r.tx.Model(&repo.WorkflowOperator{}).
		Where(&repo.WorkflowOperator{WorkflowId: workflowId, Node: node, UserId: userId}).
		Update("handled", 1).Error
}
