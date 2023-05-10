package model

import (
	"VitaTaskGo/library/config"
)

type TaskMember struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	TaskId   uint   `json:"task_id" gorm:"index:task_id"`
	UserId   uint64 `json:"user_id" gorm:"index:task_id"`
	Role     int8   `json:"role"`
	UserInfo *User  `json:"user_info,omitempty" gorm:"foreignKey:UserId"` // 定义为指针类型
}

func (receiver TaskMember) TableName() string {
	return config.Instances.Mysql.Prefix + "task_member"
}
