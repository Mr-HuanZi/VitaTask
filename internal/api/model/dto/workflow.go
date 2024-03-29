package dto

type WorkflowListQueryDto struct {
	PagingQuery
	UintId
	QueryParams
	DeletedQuery
	TypeId   []uint `json:"type_id,uint"` // 工作流类型ID
	Serials  string `json:"serials"`
	Status   string `json:"status"`
	Promoter uint64 `json:"promoter"`
	System   bool   `json:"system"`
}

type WorkflowTypeDto struct {
	ID         uint   `json:"id,omitempty"`
	Name       string `json:"name,omitempty"`
	Illustrate string `json:"illustrate,omitempty"`
	OrgId      uint   `json:"org_id,omitempty"`
	OnlyName   string `json:"only_name,omitempty"`
	System     bool   `json:"system"`
}

type WorkflowTypeQueryDto struct {
	PagingQuery
	UintId
	QueryParams
	DeletedQuery
	OnlyName string `json:"only_name,omitempty"`
}

type WorkflowTypeQueryBo struct {
	PagingQuery
	UintId
	DeletedQuery
	Name       string  `json:"name,omitempty"`
	OnlyName   string  `json:"only_name,omitempty"`
	CreateTime []int64 `json:"create_time,omitempty"`
}

type WorkflowNodeDto struct {
	UintId
	TypeId      uint   `json:"type_id,uint,omitempty"` // 工作流类型ID
	Node        int    `json:"node,int,omitempty"`     // 节点序号
	Name        string `json:"name"`
	Action      string `json:"action"`
	ActionValue string `json:"action_value"`
}

type WorkflowNodeQueryDto struct {
	UintId
	PagingQuery
	QueryParams
	DeletedQuery
	TypeId uint   `json:"type_id,omitempty" binding:"required"` // 工作流类型ID
	Action string `json:"action"`
}

type WorkflowNodeQueryBo struct {
	UintId
	PagingQuery
	DeletedQuery
	TypeId     uint    `json:"type_id,omitempty"` // 工作流类型ID
	Name       string  `json:"name,omitempty"`
	Action     string  `json:"action"`
	CreateTime []int64 `json:"create_time,omitempty"`
}

type WorkflowInitiateDto struct {
	TypeId uint   `json:"type_id,omitempty"` // 工作流类型ID
	Title  string `json:"title"`
}

type WorkflowExamineApproveDto struct {
	WorkflowId uint   `json:"workflow_id,omitempty"` // 工作流ID
	Action     string `json:"action"`                // 动作 作废 进行 驳回
}
