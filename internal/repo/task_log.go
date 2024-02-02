package repo

import (
	"VitaTaskGo/internal/api/model/dto"
)

type TaskLog struct {
	BaseModel
	ID           uint64 `gorm:"primaryKey" json:"id"`
	TaskId       uint   `json:"task_id" gorm:"index:task_id" on:"task_id"`
	OperateType  string `json:"operate_type"`
	Operator     uint64 `json:"operator"`
	OperateTime  int64  `json:"operate_time"`
	Message      string `json:"message"`
	Task         *Task  `json:"task"`
	OperatorInfo *User  `json:"operator_info" gorm:"foreignKey:ID;references:operator"`
}

func (receiver TaskLog) TableName() string {
	return GetTablePrefix() + "task_log"
}

type TaskLogRepo interface {
	Create(data *TaskLog) error
	Save(data *TaskLog) error
	Delete(id uint64) error
	Get(id uint64) (*TaskLog, error)
	UpdateField(id uint64, field string, value interface{}) error
	// DeleteByTask 清空某个任务的日志
	DeleteByTask(taskId uint) error
	PageListTaskLog(query dto.TaskLogQuery) ([]TaskLog, int64, error)
}
