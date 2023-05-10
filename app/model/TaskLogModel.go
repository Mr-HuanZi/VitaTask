package model

import (
	"VitaTaskGo/library/config"
)

type TaskLog struct {
	BaseModel
	ID           uint64 `gorm:"primaryKey" json:"id"`
	TaskId       uint   `json:"task_id" gorm:"index:task_id" on:"task_id"`
	OperateType  string `json:"operate_type"`
	Operator     uint64 `json:"operator"`
	OperateTime  int64  `json:"operate_time"`
	Message      string `json:"message"`
	Task         *Task  `json:"task"`
	OperatorInfo *User  `json:"operator_info" gorm:"foreignKey:ID;references:operator"`
}

func (receiver TaskLog) TableName() string {
	return config.Instances.Mysql.Prefix + "task_log"
}
