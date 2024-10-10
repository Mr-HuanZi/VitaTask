package repo

import (
	"VitaTaskGo/internal/api/model/dto"
	"gorm.io/gorm"
)

type WorkflowLog struct {
	BaseModel
	WorkflowId uint   `json:"workflow_id,omitempty" gorm:"index:IX_workflow_id_node"` // 工作流类型ID
	Node       int    `json:"node,omitempty" gorm:"index:IX_workflow_id_node"`        // 节点序号
	Operator   uint64 `json:"operator"`
	Nickname   string `json:"nickname"`
	Explain    string `json:"explain"`                       // 操作说明
	Action     string `json:"action" gorm:"index:IX_action"` // 节点动作
}

func (receiver *WorkflowLog) TableName() string {
	return GetTablePrefix() + "workflow_logs"
}

type WorkflowLogRepo interface {
	Create(data *WorkflowLog) error
	Save(data *WorkflowLog) error
	Delete(id uint) error
	Get(id uint) (*WorkflowLog, error)
	UpdateField(id uint, field string, value interface{}) error
	UpdateFields(id uint, values interface{}) error
	PageList(query dto.WorkflowLogQueryBo) ([]WorkflowLog, int64, error)
	// GetWorkflowAll 获取某工作流所有日志
	GetWorkflowAll(workflowId uint) ([]WorkflowLog, error)
	SetDbInstance(tx *gorm.DB)
	// GetUserHandledObj 获取用户已处理的工作流子查询对象
	GetUserHandledObj(userid uint64, action string) (tx *gorm.DB)
}
