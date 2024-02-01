package data

import (
	"VitaTaskGo/internal/repo"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TaskFilesRepo struct {
	tx  *gorm.DB
	ctx *gin.Context
}

func (r *TaskFilesRepo) Create(data *repo.TaskFiles) error {
	return r.tx.Create(&data).Error
}

func (r *TaskFilesRepo) Save(data *repo.TaskFiles) error {
	return r.tx.Save(&data).Error
}

func (r *TaskFilesRepo) Delete(id uint) error {
	return r.tx.Delete(&repo.TaskFiles{}, id).Error
}

func (r *TaskFilesRepo) Get(id uint) (*repo.TaskFiles, error) {
	var d *repo.TaskFiles
	err := r.tx.First(&d, id).Error
	return d, err
}

func (r *TaskFilesRepo) UpdateField(id uint, field string, value interface{}) error {
	return r.tx.Model(&repo.TaskFiles{}).Where("id = ?", id).Update(field, value).Error
}

func NewTaskFilesRepo(tx *gorm.DB, ctx *gin.Context) repo.TaskFilesRepo {
	return &TaskFilesRepo{
		tx:  tx,
		ctx: ctx,
	}
}
