package repo

import (
	"VitaTaskGo/internal/api/model/dto"
	"VitaTaskGo/pkg/config"
	"gorm.io/gorm"
)

type Workflow struct {
	BaseModel
	DeletedAt
	TypeId   uint   `json:"type_id,omitempty"`
	TypeName string `json:"type_name,omitempty"`
	OrgId    uint   `json:"org_id,omitempty"`
	// 工作流编号
	Serials string `json:"serials,omitempty"`
	Title   string `json:"title,omitempty"`
	// 发起人ID
	Promoter uint64 `json:"promoter,omitempty"`
	// 发起人昵称
	Nickname string `json:"nickname,omitempty"`
	// 工作流状态
	Status int `json:"status,omitempty"`
	// 当前节点
	Node int `json:"node,omitempty"`
	// 提交次数
	SubmitNum int `json:"submit_num,omitempty"`
	// 关联工作流节点表，指定用本表的Node字段关联WorkflowNode表的Node字段
	NodeInfo *WorkflowNode      `json:"node_info" gorm:"-:migration;foreignKey:Node;references:Node"`
	Operator []WorkflowOperator `json:"operator" gorm:"-:migration;WorkflowId:Node;references:ID"`
	// 状态名 英文
	StatusText string `json:"status_text,omitempty" gorm:"-"`
}

func (receiver *Workflow) TableName() string {
	return config.Instances.Mysql.Prefix + "workflow"
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
