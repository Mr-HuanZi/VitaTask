package data

import (
	"VitaTaskGo/internal/api/model/dto"
	"VitaTaskGo/internal/repo"
	"VitaTaskGo/pkg/db"
	"VitaTaskGo/pkg/exception"
	"VitaTaskGo/pkg/response"
	"VitaTaskGo/pkg/time_tool"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"time"
)

type ProjectRepo struct {
	tx  *gorm.DB
	ctx *gin.Context
}

func (r *ProjectRepo) CreateProject(data *repo.Project) error {
	return r.tx.Create(&data).Error
}

func (r *ProjectRepo) SaveProject(data *repo.Project) error {
	return r.tx.Save(&data).Error
}

func (r *ProjectRepo) DeleteProject(id uint) error {
	return r.tx.Delete(&repo.Project{}, id).Error
}

func (r *ProjectRepo) GetProject(id uint) (*repo.Project, error) {
	var project *repo.Project
	err := r.tx.First(&project, id).Error
	return project, err
}

func (r *ProjectRepo) PageListProject(dto dto.ProjectListQuery, role []uint) ([]repo.Project, int64, error) {
	var (
		projects []repo.Project = nil
		total    int64
	)

	tx := r.tx.Model(repo.Project{})

	// 查询已删除的记录
	// 回收站功能，只查询已删除的
	if dto.Deleted {
		tx = tx.Unscoped().Where("deleted_at IS NOT NULL")
	}

	// 项目id限制
	if len(role) > 0 {
		tx = tx.Where("id IN ?", role)
	}

	// 创建时间范围搜索
	if len(dto.Time) >= 2 {
		timeRange, timeRangeErr := time_tool.ParseTimeRangeToUnix(dto.Time, time.DateTime, "milli")
		if timeRangeErr != nil {
			return projects, total, exception.ErrorHandle(timeRangeErr, response.TimeParseFail)
		}
		tx = tx.Where("create_time BETWEEN ? AND ?", timeRange[0], timeRange[1])
	}
	// 项目名称模糊搜索
	if dto.Name != "" {
		tx = tx.Where("name LIKE ?", "%"+dto.Name+"%")
	}

	// 计算总记录数
	err := tx.Model(&repo.Project{}).Count(&total).Error
	if err != nil {
		return projects, 0, exception.ErrorHandle(err, response.DbQueryError)
	}

	// 查询记录
	err = tx.Model(&repo.Project{}).
		Scopes(db.Paginate(&dto.Page, &dto.PageSize)).
		Preload("Member.UserInfo").Find(&projects).Error

	return projects, total, exception.ErrorHandle(err, response.DbQueryError)
}

func (r *ProjectRepo) SimpleList(role []uint) ([]dto.ProjectSimpleList, error) {
	var simpleProjectList []dto.ProjectSimpleList

	// 如果没有提供role
	if len(role) <= 0 {
		return nil, exception.NewException(response.TooFewElements)
	}

	// 只检索当前用户所属的项目列表
	r.tx.Model(&repo.Project{}).Where("id IN ?", role).Find(&simpleProjectList)
	return simpleProjectList, nil
}

// PreloadGetProject 获取单条记录带预加载的
func (r *ProjectRepo) PreloadGetProject(id uint) (*repo.Project, error) {
	var project *repo.Project
	err := r.tx.Preload("Member.UserInfo").First(&project, id).Error
	return project, err
}

func (r *ProjectRepo) Exist(id uint) bool {
	return r.tx.Select("id").Where("id = ?", id).First(&repo.Project{}).Error == nil
}

func (r *ProjectRepo) UpdateField(id uint, field string, value interface{}) error {
	return r.tx.Model(&repo.Project{}).Where("id = ?", id).Update(field, value).Error
}

func (r *ProjectRepo) Archived(id uint) bool {
	return r.tx.Select("id").Where("id = ?", id).Where("archive = ?", 1).First(&repo.Project{}).Error == nil
}

func (r *ProjectRepo) GetUserProjects(uid uint64) ([]repo.Project, error) {
	var project []repo.Project
	err := r.tx.Model(&repo.ProjectMember{}).
		Select("Project.id,Project.name,Project.complete,Project.archive").
		Joins("Project").
		Where(&repo.ProjectMember{UserId: uid}).
		Find(&project).
		Error
	return project, err
}

func NewProjectRepo(tx *gorm.DB, ctx *gin.Context) repo.ProjectRepo {
	return &ProjectRepo{
		tx:  tx,
		ctx: ctx,
	}
}
