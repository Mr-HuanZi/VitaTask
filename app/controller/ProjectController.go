package controller

import (
	"VitaTaskGo/app/extend/user"
	"VitaTaskGo/app/logic"
	"VitaTaskGo/app/response"
	"VitaTaskGo/app/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ProjectController struct {
}

func NewProjectController() *ProjectController {
	return &ProjectController{}
}

// CreateProject 创建项目信息
func (receiver *ProjectController) CreateProject(ctx *gin.Context) {
	var post types.CreateProjectForm
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}

	ctx.JSON(http.StatusOK, response.Auto(logic.NewProjectLogic(ctx).CreateProject(post.Name, post.Leader)))
}

// EditProject 编辑项目信息
func (receiver *ProjectController) EditProject(ctx *gin.Context) {
	var post types.EditProjectForm
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}
	ctx.JSON(http.StatusOK, response.Auto(logic.NewProjectLogic(ctx).EditProject(post.ID, post.Name, post.Leader)))
}

// ProjectList 项目列表
func (receiver *ProjectController) ProjectList(ctx *gin.Context) {
	var post types.ProjectListQuery
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}
	ctx.JSON(http.StatusOK, response.Auto(logic.NewProjectLogic(ctx).GetProjectList(post, false)))
}

// SimpleList 简单项目列表
func (receiver *ProjectController) SimpleList(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, response.SuccessData(logic.NewProjectLogic(ctx).GetSimpleList()))
}

// ProjectTrash 项目回收站
func (*ProjectController) ProjectTrash(ctx *gin.Context) {
	var post types.ProjectListQuery
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}
	ctx.JSON(http.StatusOK, response.Auto(logic.NewProjectLogic(ctx).GetProjectList(post, true)))
}

// ProjectDelete 删除项目
func (receiver *ProjectController) ProjectDelete(ctx *gin.Context) {
	var post types.ProjectSingleId
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}
	if err := logic.NewProjectLogic(ctx).ProjectDelete(post.ID); err == nil {
		ctx.JSON(http.StatusOK, response.Success())
	} else {
		ctx.JSON(http.StatusOK, response.Error(err))
	}
}

// ProjectArchive 项目归档
func (receiver *ProjectController) ProjectArchive(ctx *gin.Context) {
	var post types.ProjectSingleId
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}

	ctx.JSON(http.StatusOK, response.Auto(nil, logic.NewProjectLogic(ctx).ProjectArchive(post.ID)))
}

// UnArchive 项目归档
func (receiver *ProjectController) UnArchive(ctx *gin.Context) {
	var post types.ProjectSingleId
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}

	ctx.JSON(http.StatusOK, response.Auto(nil, logic.NewProjectLogic(ctx).UnArchive(post.ID)))
}

// Star 收藏项目
func (receiver *ProjectController) Star(ctx *gin.Context) {
	var post types.ProjectSingleId
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}
	// 获取当前用户
	currUser, err := user.CurrUser(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, response.Error(err))
		return
	}
	err = logic.NewProjectMemberLogic(ctx).ProjectStar(post.ID, currUser.ID)
	if err != nil {
		ctx.JSON(http.StatusOK, response.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, response.Success())
}

// UnStart 取消收藏项目
func (receiver *ProjectController) UnStart(ctx *gin.Context) {
	var post types.ProjectSingleId
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}
	// 获取当前用户
	currUser, err := user.CurrUser(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, response.Error(err))
		return
	}
	err = logic.NewProjectMemberLogic(ctx).ProjectUnStar(post.ID, currUser.ID)
	if err != nil {
		ctx.JSON(http.StatusOK, response.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, response.Success())
}

// Transfer 移交项目
func (receiver *ProjectController) Transfer(ctx *gin.Context) {
	var post types.ProjectTransferForm
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}
	ctx.JSON(
		http.StatusOK,
		response.Auto(nil, logic.NewProjectLogic(ctx).Transfer(post.Project, post.Transferor, post.Recipient)),
	)
}

// Detail 获取项目详情
func (receiver *ProjectController) Detail(ctx *gin.Context) {
	var post types.ProjectSingleId
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}

	ctx.JSON(
		http.StatusOK,
		response.Auto(logic.NewProjectLogic(ctx).GetOneProject(post.ID)),
	)
}
