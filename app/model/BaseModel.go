package model

import "gorm.io/gorm"

type BaseModel struct {
	ID         uint  `json:"id,omitempty" gorm:"primaryKey"`
	CreateTime int64 `json:"create_time" gorm:"autoCreateTime:milli"` // 毫秒时间戳
	UpdateTime int64 `json:"-" gorm:"autoUpdateTime:milli"`           // 前端默认不显示更新时间
}

type DeletedAt struct {
	DeletedAt gorm.DeletedAt `json:"-" gorm:"size:0"` // 前端默认不显示删除时间
}
