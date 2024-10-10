package repo

import (
	"gorm.io/gorm"
)

type WorkflowOperator struct {
	ID uint `json:"id,omitempty"`
	// 操作人ID
	UserId uint64 `json:"user_id,omitempty"`
	// 操作人昵称
	Nickname string `json:"nickname,omitempty"`
	// 操作步骤
	Node int `json:"node,omitempty"`
	// 工作流ID
	WorkflowId uint `json:"workflow_id,omitempty"`
	// 是否已处理
	Handled int `json:"handled,omitempty"`
}

func (receiver *WorkflowOperator) TableName() string {
	return GetTablePrefix() + "workflow_operator"
}

type WorkflowOperatorRepo interface {
	Create(data *WorkflowOperator) error
	Save(data *WorkflowOperator) error
	Delete(id uint) error
	Get(id uint) (*WorkflowOperator, error)
	UpdateField(id uint, field string, value interface{}) error
	UpdateFields(id uint, values interface{}) error
	// GetWorkflowOperatorByNode 获取该工作流的该节点的操作人
	GetWorkflowOperatorByNode(workflowId uint, node int) ([]WorkflowOperator, error)
	SetDbInstance(tx *gorm.DB)
	// OtherOperator 该工作流的该节点是否还有除了userId以外的操作人
	OtherOperator(workflowId uint, node int, userId uint64) (bool, error)
	// RemoveWorkflowAllOperator 删除该工作流所有操作人
	RemoveWorkflowAllOperator(workflowId uint) error
	// SetHandled 将当前步骤的指定操作人改为已操作的状态
	SetHandled(workflowId uint, node int, userId uint64) error
	// GetUserTodoObj 获取用户待办列表查询对象
	GetUserTodoObj(userid uint64) *gorm.DB
}
