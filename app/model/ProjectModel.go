package model

import (
	"VitaTaskGo/library/config"
)

type Project struct {
	BaseModel
	DeletedAt
	Name     string           `json:"name,omitempty" gorm:"size:256"`
	Complete int              `json:"complete"`
	Archive  int8             `json:"archive"`
	Member   []*ProjectMember `json:"member,omitempty" gorm:"foreignKey:ProjectId"`
	Leader   *ProjectMember   `json:"leader,omitempty" gorm:"-"` // 手动获取
}

func (receiver Project) TableName() string {
	return config.Instances.Mysql.Prefix + "project"
}
