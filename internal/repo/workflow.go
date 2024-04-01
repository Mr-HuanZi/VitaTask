package repo

import (
	"VitaTaskGo/internal/api/model/dto"
	"gorm.io/gorm"
)

type Workflow struct {
	BaseModel
	DeletedAt
	TypeId   uint   `json:"type_id"`
	TypeName string `json:"type_name"`
	OrgId    uint   `json:"org_id"`
	// 工作流编号
	Serials string `json:"serials"`
	Title   string `json:"title"`
	// 发起人ID
	Promoter uint64 `json:"promoter"`
	// 发起人昵称
	Nickname string `json:"nickname"`
	// 工作流状态
	Status int `json:"status"`
	// 当前节点
	Node int `json:"node"`
	// 提交次数
	SubmitNum int `json:"submit_num"`
	// 备注 发起时写在主表的是备注，审批过程填写的叫做 说明
	Remarks string `json:"remarks"`
	// 关联工作流节点表，指定用本表的Node字段关联WorkflowNode表的Node字段
	NodeInfo *WorkflowNode      `json:"node_info" gorm:"-:migration;foreignKey:Node;references:Node"`
	Operator []WorkflowOperator `json:"operator" gorm:"-:migration;WorkflowId:Node;references:ID"`
	// 状态名 英文
	StatusText string `json:"status_text" gorm:"-"`
}

func (receiver *Workflow) TableName() string {
	return GetTablePrefix() + "workflow"
}

type WorkflowRepo interface {
	Create(data *Workflow) error
	Save(data *Workflow) error
	Delete(id uint) error
	Get(id uint) (*Workflow, error)
	UpdateField(id uint, field string, value interface{}) error
	UpdateFields(id uint, values interface{}) error
	PageList(query dto.WorkflowListQueryDto) ([]Workflow, int64, error)
	SetDbInstance(tx *gorm.DB)
	GetDayTotal(start, end int64) (int64, error)
}
