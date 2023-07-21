package biz

import (
	"VitaTaskGo/internal/pkg/config"
)

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

type DialogRepo interface {
	GetDialog(id uint) (*Dialog, error)
	CreateDialog(dialog *Dialog) error
	UpdateDialogStruct(dialog *Dialog) error
	DeleteDialog(id uint) error
	InDialog(id uint, uid uint64) bool
}
