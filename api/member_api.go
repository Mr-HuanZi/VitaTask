package api

import (
	"VitaTaskGo/internal/dto"
	"VitaTaskGo/internal/pkg/db"
	"VitaTaskGo/internal/pkg/response"
	"VitaTaskGo/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type MemberApi struct {
}

func NewMemberApi() *MemberApi {
	return &MemberApi{}
}

func (receiver MemberApi) Lists(ctx *gin.Context) {
	var (
		query dto.MemberListsQuery
	)
	if err := ctx.ShouldBindJSON(&query); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}
	ctx.JSON(http.StatusOK, response.Auto(service.NewMemberService(db.Db, ctx).Lists(query)))
}

// SimpleList 简单的成员列表
func (receiver MemberApi) SimpleList(ctx *gin.Context) {
	ctx.JSON(
		http.StatusOK,
		response.SuccessData(
			service.NewMemberService(db.Db, ctx).SimpleList(ctx.Query("key")),
		),
	)
}

// Create 创建成员
func (receiver MemberApi) Create(ctx *gin.Context) {
	var post dto.MemberCreate
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}
	ctx.JSON(http.StatusOK, response.Auto(service.NewMemberService(db.Db, ctx).Create(post)))
}

// Disable 禁用成员
func (receiver MemberApi) Disable(ctx *gin.Context) {
	var post dto.PostUid
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusOK, response.Exception(response.FormVerificationFailed))
		return
	}
	err := service.NewMemberService(db.Db, ctx).ChangeUserStatus(post.Uid, 2)
	ctx.JSON(http.StatusOK, response.Auto(nil, err))
}

// Enable 启用成员
func (receiver MemberApi) Enable(ctx *gin.Context) {
	var post dto.PostUid
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusOK, response.Exception(response.FormVerificationFailed))
		return
	}
	err := service.NewMemberService(db.Db, ctx).ChangeUserStatus(post.Uid, 1)
	ctx.JSON(http.StatusOK, response.Auto(nil, err))
}

// ResetPassword 重置用户密码
func (receiver MemberApi) ResetPassword(ctx *gin.Context) {
	var post dto.PostUid
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusOK, response.Exception(response.FormVerificationFailed))
		return
	}
	ctx.JSON(http.StatusOK, response.Auto(nil, service.NewMemberService(db.Db, ctx).ResetPassword(post.Uid)))
}

// ChangeSuper 改变一个成员的超级管理员状态
func (receiver MemberApi) ChangeSuper(ctx *gin.Context) {
	var post dto.ChangeSuperDto
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}
	// 表单校验要支持0值有点麻烦，前端给到的Super是+1的值
	ctx.JSON(http.StatusOK, response.Auto(nil, service.NewMemberService(db.Db, ctx).ChangeSuper(post.Uid, post.Super-1)))
}
