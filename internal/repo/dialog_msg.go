package repo

import (
	"VitaTaskGo/pkg/config"
)

type DialogMsg struct {
	BaseModel
	ID       uint64  `json:"id" gorm:"primaryKey"`
	DialogId uint    `json:"dialog_id" gorm:"default:0;not null;"`
	UserId   uint64  `json:"user_id" gorm:"default:0;not null"`
	Type     string  `json:"type" gorm:"size:30;default:''"`
	Content  string  `json:"content" gorm:"type:longtext"`
	Dialog   *Dialog `json:"dialog" gorm:"-:migration"`
	// 关联用户表，指定用本表的UserId字段关联User表的ID字段
	UserInfo *User `json:"user_info" gorm:"-:migration;foreignKey:ID;references:UserId"`
	DeletedAt
}

func (receiver DialogMsg) TableName() string {
	return config.Instances.Mysql.Prefix + "dialog_msg"
}

type DialogMsgRepo interface {
	GetDialogMsg(id uint) (*DialogMsg, error)
	ListDialogMsg(dialogId uint) ([]DialogMsg, error)
	CreateDialogMsg(dialogMsg *DialogMsg) error
	UpdateDialogMsgStruct(dialogMsg *DialogMsg) error
	// DeleteDialogMsg 删除对话的所有消息记录
	DeleteDialogMsg(dialogId uint) error
}
