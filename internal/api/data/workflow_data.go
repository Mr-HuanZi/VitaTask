package data

import (
	"VitaTaskGo/internal/repo"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type WorkflowDataRepo struct {
	*BaseRepo[repo.WorkflowData]
}

func NewWorkflowDataRepo(tx *gorm.DB, ctx *gin.Context) repo.WorkflowDataRepo {
	return &WorkflowDataRepo{
		BaseRepo: NewBaseRepo[repo.WorkflowData](tx, ctx),
	}
}

func (r *WorkflowDataRepo) AllData(workflowId uint) ([]repo.WorkflowData, error) {
	var data []repo.WorkflowData
	err := r.tx.Model(&repo.WorkflowData{}).Where("workflow_id = ?", workflowId).Order("create_time").Find(&data).Error
	return data, err
}
