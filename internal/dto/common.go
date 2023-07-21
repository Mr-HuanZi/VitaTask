package dto

type PagingQuery struct {
	Page     int `json:"page"`
	PageSize int `json:"pageSize" binding:"required"`
}

type PagedResult[T any] struct {
	Items []T   `json:"items"`
	Total int64 `json:"total"`
	Page  int64 `json:"page"`
}

type IntId struct {
	ID int
}

type UintId struct {
	ID uint `json:"id,omitempty"`
}

type QueryParams struct {
	Name       string   `json:"name"`
	Title      string   `json:"title"`
	CreateTime []string `json:"create_time"`
}

type DeletedQuery struct {
	Deleted bool `json:"deleted"`
}

type SingleUintRequired struct {
	ID uint `json:"id" binding:"required"`
}

// UniversalSimpleList 通用的简单列表
// 为了适配前端的下拉选择框
type UniversalSimpleList[T comparable] struct {
	Label string `json:"label"`
	Value T      `json:"value"`
}

// UniversalSimpleGroupList 通用的简单列表-带分组
// 为了适配前端的下拉选择框
type UniversalSimpleGroupList[T comparable] struct {
	Label   string                   `json:"label"`
	Options []UniversalSimpleList[T] `json:"options"`
}
