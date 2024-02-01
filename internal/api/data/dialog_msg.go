package data

import (
	"VitaTaskGo/internal/repo"
	"VitaTaskGo/pkg/exception"
	"VitaTaskGo/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DialogMsgRepo struct {
	tx  *gorm.DB
	ctx *gin.Context
}

func (receiver *DialogMsgRepo) GetDialogMsg(id uint) (*repo.DialogMsg, error) {
	var dialogMsg *repo.DialogMsg
	err := receiver.tx.Model(&repo.DialogMsg{}).First(&dialogMsg, id).Error
	return dialogMsg, exception.ErrorHandle(err, response.DialogNotExist)
}

func (receiver *DialogMsgRepo) ListDialogMsg(dialogId uint) ([]repo.DialogMsg, error) {
	var list []repo.DialogMsg
	err := receiver.tx.Model(&repo.DialogMsg{}).
		Preload("Dialog").   // 预加载
		Preload("UserInfo"). // 预加载
		Where("dialog_id = ?", dialogId).
		Find(&list).Error
	return list, err
}

func (receiver *DialogMsgRepo) CreateDialogMsg(dialogMsg *repo.DialogMsg) error {
	return receiver.tx.Create(&dialogMsg).Error
}

func (receiver *DialogMsgRepo) UpdateDialogMsgStruct(dialogMsg *repo.DialogMsg) error {
	return receiver.tx.Save(&dialogMsg).Error
}

func (receiver *DialogMsgRepo) DeleteDialogMsg(dialogId uint) error {
	return receiver.tx.Where("dialog_id = ?", dialogId).Delete(&repo.DialogMsg{}).Error
}

func NewDialogMsgRepo(tx *gorm.DB, ctx *gin.Context) repo.DialogMsgRepo {
	return &DialogMsgRepo{
		tx:  tx,
		ctx: ctx,
	}
}
