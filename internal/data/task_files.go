package data

import (
	"VitaTaskGo/internal/biz"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TaskFilesRepo struct {
	tx  *gorm.DB
	ctx *gin.Context
}

func (r *TaskFilesRepo) Create(data *biz.TaskFiles) error {
	return r.tx.Create(&data).Error
}

func (r *TaskFilesRepo) Save(data *biz.TaskFiles) error {
	return r.tx.Save(&data).Error
}

func (r *TaskFilesRepo) Delete(id uint) error {
	return r.tx.Delete(&biz.TaskFiles{}, id).Error
}

func (r *TaskFilesRepo) Get(id uint) (*biz.TaskFiles, error) {
	var d *biz.TaskFiles
	err := r.tx.First(&d, id).Error
	return d, err
}

func (r *TaskFilesRepo) UpdateField(id uint, field string, value interface{}) error {
	return r.tx.Model(&biz.TaskFiles{}).Where("id = ?", id).Update(field, value).Error
}

func NewTaskFilesRepo(tx *gorm.DB, ctx *gin.Context) biz.TaskFilesRepo {
	return &TaskFilesRepo{
		tx:  tx,
		ctx: ctx,
	}
}
