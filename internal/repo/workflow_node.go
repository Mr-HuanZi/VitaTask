package repo

import (
	"VitaTaskGo/internal/api/model/dto"
	"gorm.io/gorm"
)

type WorkflowNode struct {
	BaseModel
	DeletedAt
	TypeId      uint   `json:"type_id,omitempty" gorm:"index:type_id"` // 工作流类型ID
	Node        int    `json:"node,omitempty"`                         // 节点序号
	Name        string `json:"name"`
	Action      string `json:"action"`
	ActionValue string `json:"action_value"`
	Everyone    int    `json:"everyone"`
}

func (receiver *WorkflowNode) TableName() string {
	return GetTablePrefix() + "workflow_node"
}

type WorkflowNodeRepo interface {
	Create(data *WorkflowNode) error
	Save(data *WorkflowNode) error
	Delete(id uint) error
	Get(id uint) (*WorkflowNode, error)
	UpdateField(id uint, field string, value interface{}) error
	UpdateFields(id uint, values interface{}) error
	PageList(query dto.WorkflowNodeQueryBo) ([]WorkflowNode, int64, error)
	// GetAppointNode 获取指定节点
	GetAppointNode(typeId uint, node int) (*WorkflowNode, error)
	GetNextNode(typeId uint, currNode int) (*WorkflowNode, error)
	SetDbInstance(tx *gorm.DB)
	FirstNode(typeId uint) (*WorkflowNode, error)
}
