package repo

import (
	"VitaTaskGo/internal/api/model/dto"
	"VitaTaskGo/pkg/config"
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
	return config.Get().Mysql.Prefix + "project"
}

type ProjectRepo interface {
	CreateProject(data *Project) error
	SaveProject(data *Project) error
	DeleteProject(id uint) error
	GetProject(id uint) (*Project, error)
	PageListProject(dto dto.ProjectListQuery, role []uint) ([]Project, int64, error)
	SimpleList(role []uint) ([]dto.ProjectSimpleList, error)
	PreloadGetProject(id uint) (*Project, error)
	Exist(id uint) bool
	UpdateField(id uint, field string, value interface{}) error
	// Archived 是否归档
	Archived(id uint) bool
	GetUserProjects(uid uint64) ([]Project, error)
}
