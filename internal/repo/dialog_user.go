package repo

import (
	"VitaTaskGo/pkg/config"
)

type DialogUser struct {
	BaseModel
	DialogId uint   `json:"dialog_id" gorm:"default:0;not null;"`
	UserId   uint64 `json:"user_id" gorm:"default:0;not null"`
}

func (receiver DialogUser) TableName() string {
	return config.Instances.Mysql.Prefix + "dialog_user"
}

type DialogUserRepo interface {
	CreateDialogUser(data *DialogUser) error
	UpdateDialogUser(data *DialogUser) error
	DeleteDialogUser(dialogId uint, users []uint64) error
	GetDialogUsers(dialogId uint) ([]DialogUser, error)
	DeleteDialogAllUser(dialogId uint) error
}
