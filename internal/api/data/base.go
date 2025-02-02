package data

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// BaseRepo 泛型基础存储库，处理通用CRUD操作
type BaseRepo[T any] struct {
	tx  *gorm.DB
	ctx *gin.Context
}

func NewBaseRepo[T any](tx *gorm.DB, ctx *gin.Context) *BaseRepo[T] {
	return &BaseRepo[T]{tx: tx, ctx: ctx}
}

func (r *BaseRepo[T]) Create(data *T) error {
	return r.tx.Create(&data).Error
}

func (r *BaseRepo[T]) Save(data *T) error {
	return r.tx.Save(&data).Error
}

func (r *BaseRepo[T]) Delete(id uint) error {
	return r.tx.Delete(new(T), id).Error
}

func (r *BaseRepo[T]) Get(id uint) (*T, error) {
	var result T
	err := r.tx.First(&result, id).Error
	return &result, err
}

func (r *BaseRepo[T]) UpdateField(id uint, field string, value interface{}) error {
	return r.tx.Model(new(T)).Where("id = ?", id).Update(field, value).Error
}

func (r *BaseRepo[T]) UpdateFields(id uint, values interface{}) error {
	return r.tx.Model(new(T)).Where("id = ?", id).Updates(values).Error
}

func (r *BaseRepo[T]) SetDbInstance(tx *gorm.DB) {
	r.tx = tx
}

// FirstStringWhere 根据字符串条件查询一条数据
func (r *BaseRepo[T]) FirstStringWhere(str string, values ...interface{}) (*T, error) {
	var result T
	err := r.tx.Where(str, values...).First(&result).Error
	return &result, err
}
