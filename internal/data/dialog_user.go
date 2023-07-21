package data

import (
	"VitaTaskGo/internal/biz"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DialogUserRepo struct {
	tx  *gorm.DB
	ctx *gin.Context
}

func (receiver *DialogUserRepo) CreateDialogUser(data *biz.DialogUser) error {
	return receiver.tx.Create(&data).Error
}

func (receiver *DialogUserRepo) UpdateDialogUser(data *biz.DialogUser) error {
	return receiver.tx.Save(&data).Error
}

func (receiver *DialogUserRepo) DeleteDialogUser(dialogId uint, users []uint64) error {
	return receiver.tx.Where("dialog_id = ?", dialogId).Where("user_id IN ?", users).Delete(&biz.DialogUser{}).Error
}

func (receiver *DialogUserRepo) GetDialogUsers(dialogId uint) ([]biz.DialogUser, error) {
	var members []biz.DialogUser

	err := receiver.tx.Model(&biz.DialogUser{}).Where("dialog_id = ?", dialogId).Find(&members).Error
	return members, err
}

func (receiver *DialogUserRepo) DeleteDialogAllUser(dialogId uint) error {
	return receiver.tx.Where("dialog_id = ?", dialogId).Delete(&biz.DialogUser{}).Error
}

func NewDialogUserRepo(tx *gorm.DB, ctx *gin.Context) biz.DialogUserRepo {
	return &DialogUserRepo{
		tx:  tx,
		ctx: ctx,
	}
}
