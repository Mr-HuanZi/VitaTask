package model

import (
	"VitaTaskGo/library/config"
)

type TaskGroup struct {
	BaseModel
	DeletedAt
	ProjectId uint     `json:"project_id" gorm:"index:project_id"`
	Name      string   `json:"name" gorm:"size:256"`
	Project   *Project `json:"project,omitempty"` // 一对多（反向）
}

func (receiver TaskGroup) TableName() string {
	return config.Instances.Mysql.Prefix + "task_group"
}
