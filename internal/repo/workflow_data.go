package repo

import (
	"gorm.io/gorm"
)

type WorkflowData struct {
	BaseModel
	TypeId   uint   `json:"type_id"`
	TypeName string `json:"type_name"`
	// 工作流实例ID
	WorkflowId uint `json:"workflow_id"`
	// 所属节点ID
	NodeId uint   `json:"node_id"`
	Node   int    `json:"node"` // 节点序号
	Data   string `json:"data" gorm:"type:json"`
	Schema string `json:"schema" gorm:"type:json"`
}

func (receiver *WorkflowData) TableName() string {
	return GetTablePrefix() + "workflow_data"
}

type WorkflowDataRepo interface {
	Create(data *WorkflowData) error
	Save(data *WorkflowData) error
	Delete(id uint) error
	Get(id uint) (*WorkflowData, error)
	UpdateField(id uint, field string, value interface{}) error
	UpdateFields(id uint, values interface{}) error
	SetDbInstance(tx *gorm.DB)
	FirstStringWhere(string, ...interface{}) (*WorkflowData, error)
	AllData(workflowId uint) ([]WorkflowData, error)
}
