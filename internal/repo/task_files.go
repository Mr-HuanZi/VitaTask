package repo

import (
	"VitaTaskGo/pkg/config"
)

type TaskFiles struct {
	BaseModel
	DeletedAt
	ProjectId uint   `gorm:"index:project_id"`
	TaskId    uint   `gorm:"index:project_id"`
	UserId    uint64 `gorm:"index:project_id"`
	Filename  string `gorm:"size:256"`
	Path      string `gorm:"size:256"`
	Md5       string `gorm:"size:50"`
	Size      uint64
	Ext       string `gorm:"size:256"`
	Download  int
	Thumb     string `gorm:"size:256"`
}

func (receiver TaskFiles) TableName() string {
	return config.Instances.Mysql.Prefix + "task_files"
}

type TaskFilesRepo interface {
	Create(data *TaskFiles) error
	Save(data *TaskFiles) error
	Delete(id uint) error
	Get(id uint) (*TaskFiles, error)
	UpdateField(id uint, field string, value interface{}) error
}
