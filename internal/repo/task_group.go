package repo

import (
	"VitaTaskGo/internal/api/model/dto"
)

type TaskGroup struct {
	BaseModel
	DeletedAt
	ProjectId uint     `json:"project_id" gorm:"index:project_id"`
	Name      string   `json:"name" gorm:"size:256"`
	Project   *Project `json:"project,omitempty"` // 一对多（反向）
}

func (receiver TaskGroup) TableName() string {
	return GetTablePrefix() + "task_group"
}

type TaskGroupRepo interface {
	Create(data *TaskGroup) error
	Save(data *TaskGroup) error
	Delete(id uint) error
	Get(id uint) (*TaskGroup, error)
	UpdateField(id uint, field string, value interface{}) error
	Exist(id uint) bool
	PageListTaskLog(query dto.TaskGroupQuery) ([]TaskGroup, int64, error)
	// Detail 获取详情，带预加载的
	Detail(id uint) (*TaskGroup, error)
	SimpleList(projectId uint) ([]TaskGroup, error)
}
