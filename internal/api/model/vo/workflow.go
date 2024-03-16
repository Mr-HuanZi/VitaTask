package vo

import "VitaTaskGo/internal/repo"

type WorkflowDetailVo struct {
	Workflow     *repo.Workflow          `json:"workflow"`
	Node         *WorkflowNodeVo         `json:"node"`
	Operators    []repo.WorkflowOperator `json:"operators"`
	WorkflowType *repo.WorkflowType      `json:"workflow_type"`
}

type WorkflowNodeVo struct {
	Node        int    `json:"node,omitempty"`
	Name        string `json:"name"`
	Action      string `json:"action"`
	ActionValue string `json:"action_value"`
	Everyone    int    `json:"everyone"`
}

type WorkflowLogVo struct {
	ID         uint               `json:"id" gorm:"primaryKey"`
	CreateTime int64              `json:"create_time"`
	WorkflowId uint               `json:"workflow_id"`
	Node       int                `json:"node"`
	Operator   uint64             `json:"operator"`
	Nickname   string             `json:"nickname"`
	Explain    string             `json:"explain"`
	Action     string             `json:"action"`
	NodeInfo   *repo.WorkflowNode `json:"node_info"`
}
