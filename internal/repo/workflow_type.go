package repo

import (
	"VitaTaskGo/internal/api/model/dto"
	"VitaTaskGo/pkg/config"
	"gorm.io/gorm"
)

type WorkflowType struct {
	BaseModel
	DeletedAt
	Name       string `json:"name,omitempty"`
	Illustrate string `json:"illustrate"`
	// 所属组织。如果为空则为全局工作流
	OrgId uint `json:"org_id,omitempty"`
	// 工作流类型唯一名称。全局唯一名称，此字段用于匹配流程的模型、实例注册等，例如用作模型，表名为【flow_data_test】。该字段只需要填写【test】即可
	OnlyName string `json:"only_name,omitempty"`
	// 系统级工作流类型 1-是 0-否
	System int8 `json:"system,omitempty"`
}

func (receiver *WorkflowType) TableName() string {
	return config.Instances.Mysql.Prefix + "workflow_type"
}

type WorkflowTypeRepo interface {
	Create(data *WorkflowType) error
	Save(data *WorkflowType) error
	Delete(id uint) error
	Get(id uint) (*WorkflowType, error)
	GetByOnlyName(onlyName string) (*WorkflowType, error)
	UpdateField(id uint, field string, value interface{}) error
	UpdateFields(id uint, values interface{}) error
	PageList(query dto.WorkflowTypeQueryBo) ([]WorkflowType, int64, error)
	GetOptions(keyWords string, system bool) ([]WorkflowType, error)
	ExistByOnlyName(onlyName string) bool
	GetNotSystemIds() ([]uint, error)
	SetDbInstance(tx *gorm.DB)
}
