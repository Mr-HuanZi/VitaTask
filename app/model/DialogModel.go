package model

import "VitaTaskGo/library/config"

type Dialog struct {
	BaseModel
	Type   string `json:"type" gorm:"size:30;default:''"`
	Name   string `json:"name"`
	LastAt int64  `json:"last_at" gorm:"default:0"`
	DeletedAt
}

func (receiver Dialog) TableName() string {
	return config.Instances.Mysql.Prefix + "dialog"
}
