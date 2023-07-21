package biz

import (
	"VitaTaskGo/internal/dto"
	"VitaTaskGo/internal/pkg/config"
)

type Task struct {
	BaseModel
	DeletedAt
	ProjectId    uint          `json:"project_id" gorm:"index:project_id"`
	GroupId      uint          `json:"group_id" gorm:"index:project_id"`
	Title        string        `json:"title" gorm:"size:256"`
	Describe     string        `json:"describe,omitempty" gorm:""`
	Status       uint8         `json:"status" gorm:"index:project_id"`
	Level        uint          `json:"level" gorm:"index:project_id"`
	CompleteDate int64         `json:"complete_date" gorm:"default:null"`
	ArchivedDate int64         `json:"archived_date" gorm:"default:null"`
	StartDate    int64         `json:"start_date" gorm:"default:null"`
	EndDate      int64         `json:"end_date" gorm:"default:null"`
	EnclosureNum uint          `json:"enclosure_num"`
	DialogId     uint          `json:"dialog_id" gorm:"default:0"`
	PlanTime     []int64       `json:"plan_time" gorm:"-"`
	Project      *Project      `json:"project,omitempty"` // 一对多（反向）
	Member       []*TaskMember `json:"member,omitempty" gorm:"foreignKey:TaskId"`
	Leader       *TaskMember   `json:"leader,omitempty" gorm:"-"`       // 手动获取
	Creator      *TaskMember   `json:"creator,omitempty" gorm:"-"`      // 手动获取
	Collaborator []*TaskMember `json:"collaborator,omitempty" gorm:"-"` // 手动获取
	Group        *TaskGroup    `json:"group"`                           // 一对一
}

func (receiver Task) TableName() string {
	return config.Instances.Mysql.Prefix + "task"
}

type TaskRepo interface {
	Create(data *Task) error
	Save(data *Task) error
	Delete(id uint) error
	Get(id uint) (*Task, error)
	UpdateField(id uint, field string, value interface{}) error
	UpdateFields(id uint, values interface{}) error
	PageListProject(query dto.TaskListQueryBO) ([]Task, int64, error)
	Detail(id uint) (*Task, error)
	TaskNumber(projectId uint, status []int) (int64, error)
	GetTasksByProject(projectId uint, status []int) ([]Task, error)
	CompletedQuantity(projectId uint, completeTime []int64) (int64, error)
	CreatedQuantity(projectId uint, createTime []int64) (int64, error)
}
