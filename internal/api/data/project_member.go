package data

import (
	"VitaTaskGo/internal/api/model/dto"
	"VitaTaskGo/internal/repo"
	"VitaTaskGo/pkg/db"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type ProjectMemberRepo struct {
	tx  *gorm.DB
	ctx *gin.Context
}

func (r *ProjectMemberRepo) CreateProjectMember(data *repo.ProjectMember) error {
	return r.tx.Create(&data).Error
}

func (r *ProjectMemberRepo) SaveProjectMember(data *repo.ProjectMember) error {
	return r.tx.Save(&data).Error
}

func (r *ProjectMemberRepo) DeleteProjectMember(id uint) error {
	return r.tx.Delete(&repo.ProjectMember{}, id).Error
}

func (r *ProjectMemberRepo) GetProjectMember(projectId uint, userId uint64) (*repo.ProjectMember, error) {
	var member *repo.ProjectMember
	err := r.tx.Where("project_id = ?", projectId).Where("user_id", userId).First(&member).Error
	return member, err
}

func (r *ProjectMemberRepo) GetProjectMembers(projectId uint, userIds []uint64) ([]repo.ProjectMember, error) {
	var members []repo.ProjectMember
	err := r.tx.Model(&repo.ProjectMember{}).
		Where(&repo.ProjectMember{
			ProjectId: projectId,
		}).
		Where("user_id IN ?", userIds).
		Find(&members).Error
	return members, err
}

func (r *ProjectMemberRepo) GetProjectAllMember(projectId uint) ([]repo.ProjectMember, error) {
	var members []repo.ProjectMember
	err := r.tx.Where("project_id = ?", projectId).Find(&members).Error
	return members, err
}

func (r *ProjectMemberRepo) PageListProjectMember(query dto.ProjectMemberListQuery) ([]repo.ProjectMember, int64, error) {
	var (
		projectMembers []repo.ProjectMember
		total          int64
	)

	// 指定项目id
	tx := r.tx.Where("project_id = ?", query.ProjectId)

	// 角色搜索
	if query.Role > 0 {
		tx = tx.Where("role", query.Role)
	}

	// 用户名和用户昵称联合搜索
	if query.Username != "" || query.Nickname != "" {
		// Joins连接
		tx = tx.Joins("UserInfo")

		if query.Username != "" {
			tx = tx.Clauses(clause.Like{Column: "UserInfo.user_login", Value: "%" + query.Username + "%"})
		}
		if query.Nickname != "" {
			tx = tx.Clauses(clause.Like{Column: "UserInfo.user_nickname", Value: "%" + query.Nickname + "%"})
		}
	} else {
		// 预加载
		tx = tx.Preload("UserInfo")
	}

	// 计算总记录数
	err := tx.Model(&repo.ProjectMember{}).Count(&total).Error

	if err == nil {
		err = tx.Model(&repo.ProjectMember{}).Preload("Project").Scopes(db.Paginate(&query.Page, &query.PageSize)).Order("role ASC").Find(&projectMembers).Error
	}

	return projectMembers, total, err
}

func (r *ProjectMemberRepo) UpdateField(id uint, field string, value interface{}) error {
	return r.tx.Model(&repo.ProjectMember{}).Where("id = ?", id).Update(field, value).Error
}

func (r *ProjectMemberRepo) InProject(projectId uint, userId uint64, roles []int) bool {
	tx := r.tx.Select("id").Where(&repo.ProjectMember{ProjectId: projectId, UserId: userId})
	if len(roles) > 0 {
		tx = tx.Where("role IN ?", roles)
	}

	// 有记录就说明查到了
	return tx.First(&repo.ProjectMember{}).Error == nil
}

func (r *ProjectMemberRepo) GetMembersByRole(projectId uint, roles []int) ([]repo.ProjectMember, error) {
	var members []repo.ProjectMember
	tx := r.tx.Model(&repo.ProjectMember{}).Where(&repo.ProjectMember{ProjectId: projectId})
	if len(roles) > 0 {
		tx = tx.Where("role IN ?", roles)
	}

	err := tx.Find(&members).Error
	return members, err
}

func NewProjectMemberRepo(tx *gorm.DB, ctx *gin.Context) repo.ProjectMemberRepo {
	return &ProjectMemberRepo{
		tx:  tx,
		ctx: ctx,
	}
}
