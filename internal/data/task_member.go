package data

import (
	"VitaTaskGo/internal/biz"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TaskMemberRepo struct {
	tx  *gorm.DB
	ctx *gin.Context
}

func (r *TaskMemberRepo) Create(data *biz.TaskMember) error {
	return r.tx.Create(&data).Error
}

func (r *TaskMemberRepo) Save(data *biz.TaskMember) error {
	return r.tx.Save(&data).Error
}

func (r *TaskMemberRepo) Delete(id uint) error {
	return r.tx.Delete(&biz.TaskMember{}, id).Error
}

func (r *TaskMemberRepo) Get(id uint) (*biz.TaskMember, error) {
	var d *biz.TaskMember
	err := r.tx.First(&d, id).Error
	return d, err
}

func (r *TaskMemberRepo) UpdateField(id uint, field string, value interface{}) error {
	return r.tx.Model(&biz.TaskMember{}).Where("id = ?", id).Update(field, value).Error
}

func (r *TaskMemberRepo) GetTaskMember(taskId uint, userId uint64) (*biz.TaskMember, error) {
	var member *biz.TaskMember
	err := r.tx.Where(&biz.TaskMember{TaskId: taskId, UserId: userId}).First(&member).Error

	return member, err
}

func (r *TaskMemberRepo) GetTaskMembers(taskId uint, userIds []uint64) ([]biz.TaskMember, error) {
	var members []biz.TaskMember
	err := r.tx.Model(&biz.TaskMember{}).
		Where(&biz.TaskMember{
			TaskId: taskId,
		}).
		Where("user_id IN ?", userIds).
		Find(&members).Error

	return members, err
}

func (r *TaskMemberRepo) GetTaskAllMember(taskId uint) ([]biz.TaskMember, error) {
	var members []biz.TaskMember
	err := r.tx.Where(&biz.TaskMember{
		TaskId: taskId,
	}).Find(&members).Error

	return members, err
}

func (r *TaskMemberRepo) InTask(taskId uint, userId uint64, roles []int) bool {
	tx := r.tx.Select("id").Where(&biz.TaskMember{TaskId: taskId, UserId: userId})
	if len(roles) > 0 {
		tx = tx.Where("role IN ?", roles)
	}

	// 有记录就说明查到了
	return tx.First(&biz.TaskMember{}).Error == nil
}

func (r *TaskMemberRepo) GetMembersByRole(taskId uint, roles []int) ([]biz.TaskMember, error) {
	var members []biz.TaskMember
	tx := r.tx.Model(&biz.TaskMember{}).Where(&biz.TaskMember{TaskId: taskId})
	if len(roles) > 0 {
		tx = tx.Where("role IN ?", roles)
	}

	err := tx.Find(&members).Error
	return members, err
}

func (r *TaskMemberRepo) GetTaskIdsByUsers(userIds []uint64, role []int) ([]uint, error) {
	var taskIds []uint
	err := r.tx.Model(&biz.TaskMember{}).
		Select("task_id").
		Where("user_id IN ?", userIds).
		Where("role IN ?", role).
		Find(&taskIds).Error
	return taskIds, err
}

func NewTaskMemberRepo(tx *gorm.DB, ctx *gin.Context) biz.TaskMemberRepo {
	return &TaskMemberRepo{
		tx:  tx,
		ctx: ctx,
	}
}
