package data

import (
	"VitaTaskGo/internal/biz"
	"VitaTaskGo/internal/pkg/exception"
	"VitaTaskGo/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DialogMsgRepo struct {
	tx  *gorm.DB
	ctx *gin.Context
}

func (receiver *DialogMsgRepo) GetDialogMsg(id uint) (*biz.DialogMsg, error) {
	var dialogMsg *biz.DialogMsg
	err := receiver.tx.Model(&biz.DialogMsg{}).First(&dialogMsg, id).Error
	return dialogMsg, exception.ErrorHandle(err, response.DialogNotExist)
}

func (receiver *DialogMsgRepo) ListDialogMsg(dialogId uint) ([]biz.DialogMsg, error) {
	var list []biz.DialogMsg
	err := receiver.tx.Model(&biz.DialogMsg{}).
		Preload("Dialog").   // 预加载
		Preload("UserInfo"). // 预加载
		Where("dialog_id = ?", dialogId).
		Find(&list).Error
	return list, err
}

func (receiver *DialogMsgRepo) CreateDialogMsg(dialogMsg *biz.DialogMsg) error {
	return receiver.tx.Create(&dialogMsg).Error
}

func (receiver *DialogMsgRepo) UpdateDialogMsgStruct(dialogMsg *biz.DialogMsg) error {
	return receiver.tx.Save(&dialogMsg).Error
}

func (receiver *DialogMsgRepo) DeleteDialogMsg(dialogId uint) error {
	return receiver.tx.Where("dialog_id = ?", dialogId).Delete(&biz.DialogMsg{}).Error
}

func NewDialogMsgRepo(tx *gorm.DB, ctx *gin.Context) biz.DialogMsgRepo {
	return &DialogMsgRepo{
		tx:  tx,
		ctx: ctx,
	}
}
