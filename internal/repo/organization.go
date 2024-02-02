package repo

import (
	"time"
)

type Organization struct {
	BaseModel
	Name         string `gorm:"size:256"`
	ParentId     int    `gorm:"index:parent_id"`
	Type         int8   `gorm:"index:parent_id"`
	Addr         string `gorm:"size:256"`
	RegisterDate time.Time
}

func (receiver Organization) TableName() string {
	return GetTablePrefix() + "organization"
}
