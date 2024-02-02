package repo

import (
	"VitaTaskGo/pkg/config"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID         uint  `json:"id,omitempty" gorm:"primaryKey"`
	CreateTime int64 `json:"create_time" gorm:"autoCreateTime:milli"` // 毫秒时间戳
	UpdateTime int64 `json:"-" gorm:"autoUpdateTime:milli"`           // 前端默认不显示更新时间
}

type DeletedAt struct {
	DeletedAt gorm.DeletedAt `json:"-" gorm:"size:0"` // 前端默认不显示删除时间
}

// BaseRepo 基础Repo
// 泛型T表示 模型具体类型，如Task、Project
// 泛型E表示 主键数据类型，如int、uint
// todo 学艺不精，暂时不知道怎么做到通用
type BaseRepo[T any, E comparable] interface {
	Create(data *T) error
	Save(data *T) error
	Delete(id E) error
	Get(id E) (*T, error)
	UpdateField(id E, field string, value interface{}) error
}

// GetTablePrefix 获取数据库表前缀
func GetTablePrefix() string {
	return config.Get().Mysql.Prefix
}
