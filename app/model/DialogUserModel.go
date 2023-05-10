package model

import "VitaTaskGo/library/config"

type DialogUser struct {
	BaseModel
	DialogId uint   `json:"dialog_id" gorm:"default:0;not null;"`
	UserId   uint64 `json:"user_id" gorm:"default:0;not null"`
}

func (receiver DialogUser) TableName() string {
	return config.Instances.Mysql.Prefix + "dialog_user"
}
