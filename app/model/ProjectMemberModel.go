package model

import (
	"VitaTaskGo/library/config"
)

type ProjectMember struct {
	ID        uint   `json:"id,omitempty" gorm:"primaryKey"`
	ProjectId uint   `json:"projectId,omitempty" gorm:"index:project_id"`
	UserId    uint64 `json:"userId,omitempty" gorm:"index:project_id"`
	Role      int8   `json:"role,omitempty"`
	UserInfo  *User  `json:"userInfo,omitempty" gorm:"foreignKey:UserId"` // 定义为指针类型
}

func (receiver ProjectMember) TableName() string {
	return config.Instances.Mysql.Prefix + "project_member"
}
