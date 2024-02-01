package data

import (
	"VitaTaskGo/internal/repo"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DialogUserRepo struct {
	tx  *gorm.DB
	ctx *gin.Context
}

func (receiver *DialogUserRepo) CreateDialogUser(data *repo.DialogUser) error {
	return receiver.tx.Create(&data).Error
}

func (receiver *DialogUserRepo) UpdateDialogUser(data *repo.DialogUser) error {
	return receiver.tx.Save(&data).Error
}

func (receiver *DialogUserRepo) DeleteDialogUser(dialogId uint, users []uint64) error {
	return receiver.tx.Where("dialog_id = ?", dialogId).Where("user_id IN ?", users).Delete(&repo.DialogUser{}).Error
}

func (receiver *DialogUserRepo) GetDialogUsers(dialogId uint) ([]repo.DialogUser, error) {
	var members []repo.DialogUser

	err := receiver.tx.Model(&repo.DialogUser{}).Where("dialog_id = ?", dialogId).Find(&members).Error
	return members, err
}

func (receiver *DialogUserRepo) DeleteDialogAllUser(dialogId uint) error {
	return receiver.tx.Where("dialog_id = ?", dialogId).Delete(&repo.DialogUser{}).Error
}

func NewDialogUserRepo(tx *gorm.DB, ctx *gin.Context) repo.DialogUserRepo {
	return &DialogUserRepo{
		tx:  tx,
		ctx: ctx,
	}
}
