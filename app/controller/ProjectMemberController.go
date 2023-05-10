package controller

import (
	"VitaTaskGo/app/logic"
	"VitaTaskGo/app/response"
	"VitaTaskGo/app/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ProjectMemberController struct {
}

func NewProjectMemberController() *ProjectMemberController {
	return &ProjectMemberController{}
}

func (receiver ProjectMemberController) List(ctx *gin.Context) {
	var post types.ProjectMemberListQuery
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}
	list := logic.NewProjectMemberLogic(ctx).GetMembers(post)
	ctx.JSON(http.StatusOK, response.SuccessData(list))
}

func (receiver ProjectMemberController) Remove(ctx *gin.Context) {
	var post types.ProjectMemberBind
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}

	ctx.JSON(http.StatusOK, response.Auto(nil, logic.NewProjectMemberLogic(ctx).Remove(post.ProjectId, post.UserId, post.Role)))
}

// Bind 新增普通成员
func (receiver ProjectMemberController) Bind(ctx *gin.Context) {
	var post types.ProjectMemberBind
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}
	ctx.JSON(http.StatusOK, response.Auto(nil, logic.NewProjectMemberLogic(ctx).Bind(post.ProjectId, post.UserId, post.Role)))
}
