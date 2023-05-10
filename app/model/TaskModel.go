package model

import (
	"VitaTaskGo/library/config"
)

type Task struct {
	BaseModel
	DeletedAt
	ID           uint          `json:"id,omitempty" gorm:"primaryKey"`
	ProjectId    uint          `json:"project_id" gorm:"index:project_id"`
	GroupId      uint          `json:"group_id" gorm:"index:project_id"`
	Title        string        `json:"title" gorm:"size:256"`
	Describe     string        `json:"describe,omitempty"`
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
