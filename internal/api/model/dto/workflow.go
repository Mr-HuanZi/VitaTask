package dto

type WorkflowListQueryDto struct {
	PagingQuery
	UintId
	QueryParams
	DeletedQuery
	TypeId   []uint `json:"type_id,uint"` // 工作流类型ID
	Serials  string `json:"serials"`
	Status   int    `json:"status"`
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
	// 流转模式 1-顺序流转 2-自由流转
	CirculationMode int8 `json:"circulation_mode,omitempty" binding:"required"`
}

type WorkflowTypeQueryDto struct {
	PagingQuery
	UintId
	QueryParams
	DeletedQuery
	System    bool    `json:"system"`
	TimeRange []int64 `json:"time_range"`
}

type WorkflowNodeDto struct {
	UintId
	TypeId      uint   `json:"type_id,uint,omitempty" binding:"required"`  // 工作流类型ID
	Node        int    `json:"node,int,omitempty" binding:"min=0,max=999"` // 节点序号
	Name        string `json:"name" binding:"required"`
	Action      string `json:"action"`
	ActionValue string `json:"action_value"`
	End         uint8  `json:"end" binding:"oneof=0 1"` // 是否为结束节点 1-是
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

type WorkflowNodeSaveFormDto struct {
	UintId
	Schema string `json:"schema"`
}

type WorkflowInitiateDto struct {
	TypeId   uint        `json:"type_id,omitempty"` // 工作流类型ID
	Title    string      `json:"title"`
	Data     interface{} `json:"data"` // 数据
	Remarks  string      `json:"remarks"`
	MoreData interface{} `json:"more_data"` // 额外的表单数据
}

type WorkflowExamineApproveDto struct {
	Id       uint        `json:"id"`        // 工作流ID
	Action   string      `json:"action"`    // 动作 作废 进行 驳回
	Explain  string      `json:"explain"`   // 说明
	Node     int         `json:"node"`      // 退回到哪个节点
	MoreData interface{} `json:"more_data"` // 额外的表单数据
}

type WorkflowLogQueryDto struct {
	UintId
	PagingQuery
	QueryParams
	WorkflowId uint   `json:"workflow_id,omitempty"` // 工作流ID
	Node       int    `json:"node,omitempty"`
	Action     string `json:"action"`
	Operator   uint64 `json:"operator,omitempty"`
}

type WorkflowLogQueryBo struct {
	UintId
	PagingQuery
	WorkflowId uint     `json:"workflow_id,omitempty"` // 工作流ID
	Node       int      `json:"node,omitempty"`
	Action     string   `json:"action"`
	Operator   []uint64 `json:"operator,omitempty"`
	CreateTime []int64  `json:"create_time,omitempty"`
}

type WorkflowNodeSaveCirculationDto struct {
	TypeId      uint            `json:"type_id,omitempty" binding:"required"` // 工作流类型ID
	Circulation map[uint][]uint `json:"circulation" binding:"required,gt=0"`  // 节点流转配置 key:节点ID value: []节点ID切片
}
