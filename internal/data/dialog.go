package data

import (
	"VitaTaskGo/internal/biz"
	"VitaTaskGo/internal/pkg/exception"
	"VitaTaskGo/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DialogRepo struct {
	tx  *gorm.DB
	ctx *gin.Context
}

func (receiver *DialogRepo) GetDialog(id uint) (*biz.Dialog, error) {
	var dialog *biz.Dialog
	err := receiver.tx.Model(&biz.Dialog{}).First(&dialog, id).Error
	return dialog, exception.ErrorHandle(err, response.DialogNotExist)
}

func (receiver *DialogRepo) CreateDialog(dialog *biz.Dialog) error {
	return receiver.tx.Create(&dialog).Error
}

func (receiver *DialogRepo) UpdateDialogStruct(dialog *biz.Dialog) error {
	return receiver.tx.Save(&dialog).Error
}

func (receiver *DialogRepo) DeleteDialog(id uint) error {
	return receiver.tx.Delete(&biz.Dialog{}, id).Error
}

func (receiver *DialogRepo) InDialog(id uint, uid uint64) bool {
	//var count int64
	//receiver.tx.Model(&biz.DialogUser{}).
	//	Where("dialog_id = ?", id).
	//	Where("user_id = ?", uid).
	//	Count(&count)
	return receiver.tx.Select("id").
		Where("dialog_id = ?", id).
		Where("user_id = ?", uid).
		First(&biz.DialogUser{}).
		Error == nil
	//return count > 0
}

func NewDialogRepo(tx *gorm.DB, ctx *gin.Context) biz.DialogRepo {
	return &DialogRepo{
		tx:  tx,
		ctx: ctx,
	}
}
