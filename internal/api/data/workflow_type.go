package data

import (
	"VitaTaskGo/internal/api/model/dto"
	"VitaTaskGo/internal/repo"
	"VitaTaskGo/pkg/db"
	"VitaTaskGo/pkg/exception"
	"VitaTaskGo/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type WorkflowTypeRepo struct {
	tx  *gorm.DB
	ctx *gin.Context
}

func NewWorkflowTypeRepo(tx *gorm.DB, ctx *gin.Context) repo.WorkflowTypeRepo {
	return &WorkflowTypeRepo{
		tx:  tx,
		ctx: ctx,
	}
}

func (r *WorkflowTypeRepo) Create(data *repo.WorkflowType) error {
	return r.tx.Create(&data).Error
}

func (r *WorkflowTypeRepo) Save(data *repo.WorkflowType) error {
	return r.tx.Save(&data).Error
}

func (r *WorkflowTypeRepo) Delete(id uint) error {
	return r.tx.Delete(&repo.WorkflowType{}, id).Error
}

func (r *WorkflowTypeRepo) Get(id uint) (*repo.WorkflowType, error) {
	var d *repo.WorkflowType
	err := r.tx.First(&d, id).Error
	return d, err
}

func (r *WorkflowTypeRepo) GetByOnlyName(onlyName string) (*repo.WorkflowType, error) {
	var d *repo.WorkflowType
	err := r.tx.Where(&repo.WorkflowType{OnlyName: onlyName}).First(&d).Error
	return d, err
}

func (r *WorkflowTypeRepo) UpdateField(id uint, field string, value interface{}) error {
	return r.tx.Model(&repo.WorkflowType{}).Where("id = ?", id).Update(field, value).Error
}

func (r *WorkflowTypeRepo) UpdateFields(id uint, values interface{}) error {
	return r.tx.Model(&repo.WorkflowType{}).Where("id = ?", id).Updates(values).Error
}

func (r *WorkflowTypeRepo) PageList(query dto.WorkflowTypeQueryBo) ([]repo.WorkflowType, int64, error) {
	var (
		list  []repo.WorkflowType = nil
		total int64
	)

	tx := r.tx.Model(repo.WorkflowType{})

	if query.ID > 0 {
		tx = tx.Where("id = ?", query.ID)
	}

	if len(query.CreateTime) >= 2 {
		tx = tx.Where(
			"create_time BETWEEN ? AND ?",
			query.CreateTime[0],
			query.CreateTime[1],
		)
	}

	if len(query.Name) > 0 {
		tx = tx.Where("name LIKE ?", "%"+query.Name+"%")
	}

	if len(query.OnlyName) > 0 {
		tx = tx.Where("only_name LIKE ?", "%"+query.OnlyName+"%")
	}

	// 查询已删除的记录
	if query.Deleted {
		tx = tx.Unscoped().Where("deleted_at IS NOT NULL")
	}

	// 计算总记录数
	err := tx.Count(&total).Error
	if err != nil {
		return list, 0, exception.ErrorHandle(err, response.DbQueryError)
	}

	// 查询记录
	err = tx.Scopes(db.Paginate(&query.Page, &query.PageSize)).
		Order("system ASC").
		Order("create_time DESC").
		Find(&list).Error

	return list, total, exception.ErrorHandle(err, response.DbQueryError)
}

func (r *WorkflowTypeRepo) GetOptions(keyWords string, system bool) ([]repo.WorkflowType, error) {
	var (
		list []repo.WorkflowType = nil
	)

	tx := r.tx.Model(repo.WorkflowType{})
	if len(keyWords) > 0 {
		tx = tx.Where("name LIKE ?", "%"+keyWords+"%")
	}
	if !system {
		// 只搜索非系统内置
		tx.Where("system", 0)
	}

	err := tx.Select("id", "name").Order("system ASC").Order("create_time DESC").Find(&list).Error
	return list, err
}

func (r *WorkflowTypeRepo) GetNotSystemIds() ([]uint, error) {
	var ids []uint
	err := r.tx.Model(&repo.WorkflowType{}).Where("system", 0).Pluck("id", &ids).Error
	return ids, err
}

func (r *WorkflowTypeRepo) ExistByOnlyName(onlyName string) bool {
	// 有记录就说明查到了
	return r.tx.Select("id").Where(&repo.WorkflowType{OnlyName: onlyName}).First(&repo.WorkflowType{}).Error == nil
}

func (r *WorkflowTypeRepo) SetDbInstance(tx *gorm.DB) {
	r.tx = tx
}
