package repo

import (
	"VitaTaskGo/internal/api/model/dto"
)

type ProjectMember struct {
	ID        uint   `json:"id,omitempty" gorm:"primaryKey"`
	ProjectId uint   `json:"projectId,omitempty" gorm:"index:project_id"`
	UserId    uint64 `json:"userId,omitempty" gorm:"index:project_id"`
	Role      int8   `json:"role,omitempty"`
	// -:migration 在迁移时忽略该字段
	UserInfo *User `json:"userInfo,omitempty" gorm:"-:migration;foreignKey:ID;foreignKey:UserId"` // 定义为指针类型
	// 关联项目表，指定用本表的ProjectId字段关联Project表的ID字段
	Project *Project `json:"-" gorm:"-:migration;foreignKey:ID;references:ProjectId"`
}

func (receiver ProjectMember) TableName() string {
	return GetTablePrefix() + "project_member"
}

type ProjectMemberRepo interface {
	CreateProjectMember(data *ProjectMember) error
	SaveProjectMember(data *ProjectMember) error
	DeleteProjectMember(id uint) error
	GetProjectMember(p uint, u uint64) (*ProjectMember, error)
	GetProjectAllMember(projectId uint) ([]ProjectMember, error)
	PageListProjectMember(query dto.ProjectMemberListQuery) ([]ProjectMember, int64, error)
	GetProjectMembers(projectId uint, userIds []uint64) ([]ProjectMember, error)
	UpdateField(id uint, field string, value interface{}) error
	InProject(projectId uint, userId uint64, roles []int) bool
	GetMembersByRole(projectId uint, roles []int) ([]ProjectMember, error)
}
