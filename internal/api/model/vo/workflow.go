package vo

import "VitaTaskGo/internal/repo"

type WorkflowDetailVo struct {
	Workflow     *repo.Workflow          `json:"workflow"`
	Node         *WorkflowNodeVo         `json:"node"`
	Operators    []repo.WorkflowOperator `json:"operators"`
	WorkflowType *repo.WorkflowType      `json:"workflow_type"`
}

type WorkflowNodeVo struct {
	ID           uint                `json:"id,omitempty" gorm:"primaryKey"`
	Node         int                 `json:"node,omitempty"`
	Name         string              `json:"name"`
	Action       string              `json:"action"`
	ActionValue  string              `json:"action_value"`
	ActionOption *OptionItem[string] `json:"action_option"`
	Everyone     int                 `json:"everyone"`
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

type WorkflowFootprintOperatorVo struct {
	Uid      uint64 `json:"uid"`
	Nickname string `json:"nickname"`
}

type WorkflowFootprintVo struct {
	Node      int                           `json:"node"`
	Name      string                        `json:"name"`
	Curr      bool                          `json:"curr"` // 是否当前节点
	Operators []WorkflowFootprintOperatorVo `json:"operators"`
	Explain   string                        `json:"explain"`
	Time      int64                         `json:"time"`
}
