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

type WorkflowListVo struct {
	ID         uint   `json:"id,omitempty"`
	CreateTime int64  `json:"create_time"` // 毫秒时间戳
	Name       string `json:"name,omitempty"`
	// 描述
	Illustrate string `json:"illustrate"`
	// 所属组织。如果为空则为全局工作流
	OrgId uint `json:"org_id,omitempty"`
	// 工作流类型唯一名称。全局唯一名称，此字段用于匹配流程的模型、实例注册等，例如用作模型，表名为【flow_data_test】。该字段只需要填写【test】即可
	OnlyName string `json:"only_name,omitempty"`
	// 系统内置工作流类型(不允许前端修改) 1-是 0-否
	System       int8  `json:"system,omitempty"`
	UsedQuantity int64 `json:"used_quantity"`
}
