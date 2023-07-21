package data

import (
	"VitaTaskGo/internal/biz"
	"VitaTaskGo/internal/dto"
	"VitaTaskGo/internal/pkg/db"
	"VitaTaskGo/internal/pkg/exception"
	"VitaTaskGo/internal/pkg/response"
	"VitaTaskGo/internal/pkg/time_tool"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"time"
)

type ProjectRepo struct {
	tx  *gorm.DB
	ctx *gin.Context
}

func (r *ProjectRepo) CreateProject(data *biz.Project) error {
	return r.tx.Create(&data).Error
}

func (r *ProjectRepo) SaveProject(data *biz.Project) error {
	return r.tx.Save(&data).Error
}

func (r *ProjectRepo) DeleteProject(id uint) error {
	return r.tx.Delete(&biz.Project{}, id).Error
}

func (r *ProjectRepo) GetProject(id uint) (*biz.Project, error) {
	var project *biz.Project
	err := r.tx.First(&project, id).Error
	return project, err
}

func (r *ProjectRepo) PageListProject(dto dto.ProjectListQuery, role []uint) ([]biz.Project, int64, error) {
	var (
		projects []biz.Project = nil
		total    int64
	)

	tx := r.tx.Model(biz.Project{})

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
	err := tx.Model(&biz.Project{}).Count(&total).Error
	if err != nil {
		return projects, 0, exception.ErrorHandle(err, response.DbQueryError)
	}

	// 查询记录
	err = tx.Model(&biz.Project{}).
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
	r.tx.Model(&biz.Project{}).Where("id IN ?", role).Find(&simpleProjectList)
	return simpleProjectList, nil
}

// PreloadGetProject 获取单条记录带预加载的
func (r *ProjectRepo) PreloadGetProject(id uint) (*biz.Project, error) {
	var project *biz.Project
	err := r.tx.Preload("Member.UserInfo").First(&project, id).Error
	return project, err
}

func (r *ProjectRepo) Exist(id uint) bool {
	return r.tx.Select("id").Where("id = ?", id).First(&biz.Project{}).Error == nil
}

func (r *ProjectRepo) UpdateField(id uint, field string, value interface{}) error {
	return r.tx.Model(&biz.Project{}).Where("id = ?", id).Update(field, value).Error
}

func (r *ProjectRepo) Archived(id uint) bool {
	return r.tx.Select("id").Where("id = ?", id).Where("archive = ?", 1).First(&biz.Project{}).Error == nil
}

func (r *ProjectRepo) GetUserProjects(uid uint64) ([]biz.Project, error) {
	var project []biz.Project
	err := r.tx.Model(&biz.ProjectMember{}).Joins("Project").Where(&biz.ProjectMember{UserId: uid}).Find(&project).Error
	return project, err
}

func NewProjectRepo(tx *gorm.DB, ctx *gin.Context) biz.ProjectRepo {
	return &ProjectRepo{
		tx:  tx,
		ctx: ctx,
	}
}
