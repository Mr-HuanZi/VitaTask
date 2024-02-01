package data

import (
	"VitaTaskGo/internal/repo"
	"VitaTaskGo/pkg/exception"
	"VitaTaskGo/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DialogRepo struct {
	tx  *gorm.DB
	ctx *gin.Context
}

func (receiver *DialogRepo) GetDialog(id uint) (*repo.Dialog, error) {
	var dialog *repo.Dialog
	err := receiver.tx.Model(&repo.Dialog{}).First(&dialog, id).Error
	return dialog, exception.ErrorHandle(err, response.DialogNotExist)
}

func (receiver *DialogRepo) CreateDialog(dialog *repo.Dialog) error {
	return receiver.tx.Create(&dialog).Error
}

func (receiver *DialogRepo) UpdateDialogStruct(dialog *repo.Dialog) error {
	return receiver.tx.Save(&dialog).Error
}

func (receiver *DialogRepo) DeleteDialog(id uint) error {
	return receiver.tx.Delete(&repo.Dialog{}, id).Error
}

func (receiver *DialogRepo) InDialog(id uint, uid uint64) bool {
	//var count int64
	//receiver.tx.Model(&repo.DialogUser{}).
	//	Where("dialog_id = ?", id).
	//	Where("user_id = ?", uid).
	//	Count(&count)
	return receiver.tx.Select("id").
		Where("dialog_id = ?", id).
		Where("user_id = ?", uid).
		First(&repo.DialogUser{}).
		Error == nil
	//return count > 0
}

func NewDialogRepo(tx *gorm.DB, ctx *gin.Context) repo.DialogRepo {
	return &DialogRepo{
		tx:  tx,
		ctx: ctx,
	}
}
