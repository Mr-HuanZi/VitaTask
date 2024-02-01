package repo

import (
	"VitaTaskGo/pkg/config"
)

type TaskMember struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	TaskId   uint   `json:"task_id" gorm:"index:task_id"`
	UserId   uint64 `json:"user_id" gorm:"index:task_id"`
	Role     int8   `json:"role"`
	UserInfo *User  `json:"user_info,omitempty" gorm:"foreignKey:UserId"` // 定义为指针类型
}

func (receiver TaskMember) TableName() string {
	return config.Instances.Mysql.Prefix + "task_member"
}

type TaskMemberRepo interface {
	Create(data *TaskMember) error
	Save(data *TaskMember) error
	Delete(id uint) error
	Get(id uint) (*TaskMember, error)
	UpdateField(id uint, field string, value interface{}) error
	GetTaskMember(taskId uint, userId uint64) (*TaskMember, error)
	GetTaskMembers(taskId uint, userIds []uint64) ([]TaskMember, error)
	GetTaskAllMember(taskId uint) ([]TaskMember, error)
	InTask(taskId uint, userId uint64, roles []int) bool
	GetMembersByRole(taskId uint, roles []int) ([]TaskMember, error)
	GetTaskIdsByUsers(userIds []uint64, role []int) ([]uint, error)
}
