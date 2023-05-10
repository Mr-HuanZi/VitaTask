package controller

import (
	"VitaTaskGo/app/logic"
	"VitaTaskGo/app/response"
	"VitaTaskGo/app/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

type MemberController struct {
}

func NewMemberController() *MemberController {
	return &MemberController{}
}

func (receiver MemberController) Lists(ctx *gin.Context) {
	var (
		query types.MemberListsQuery
	)
	if err := ctx.ShouldBindJSON(&query); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}
	ctx.JSON(http.StatusOK, response.Auto(logic.NewMemberLogic(ctx).Lists(query)))
}

// SimpleList 简单的成员列表
func (receiver MemberController) SimpleList(ctx *gin.Context) {
	ctx.JSON(
		http.StatusOK,
		response.SuccessData(
			logic.NewMemberLogic(ctx).SimpleList(ctx.Query("key")),
		),
	)
}

// Create 创建成员
func (receiver MemberController) Create(ctx *gin.Context) {
	var post types.MemberCreate
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}
	ctx.JSON(http.StatusOK, response.Auto(logic.NewMemberLogic(ctx).Create(post)))
}

// Disable 禁用成员
func (receiver MemberController) Disable(ctx *gin.Context) {
	var post types.PostUid
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusOK, response.Exception(response.FormVerificationFailed))
		return
	}
	err := logic.NewMemberLogic(ctx).ChangeUserStatus(post.Uid, 2)
	ctx.JSON(http.StatusOK, response.Auto(nil, err))
}

// Enable 启用成员
func (receiver MemberController) Enable(ctx *gin.Context) {
	var post types.PostUid
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusOK, response.Exception(response.FormVerificationFailed))
		return
	}
	err := logic.NewMemberLogic(ctx).ChangeUserStatus(post.Uid, 1)
	ctx.JSON(http.StatusOK, response.Auto(nil, err))
}

// ResetPassword 重置用户密码
func (receiver MemberController) ResetPassword(ctx *gin.Context) {
	var post types.PostUid
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusOK, response.Exception(response.FormVerificationFailed))
		return
	}
	ctx.JSON(http.StatusOK, response.Auto(nil, logic.NewMemberLogic(ctx).ResetPassword(post.Uid)))
}

// ChangeSuper 改变一个成员的超级管理员状态
func (receiver MemberController) ChangeSuper(ctx *gin.Context) {
	var post types.ChangeSuperDto
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}
	// 表单校验要支持0值有点麻烦，前端给到的Super是+1的值
	ctx.JSON(http.StatusOK, response.Auto(nil, logic.NewMemberLogic(ctx).ChangeSuper(post.Uid, post.Super-1)))
}
