package api

import (
	"VitaTaskGo/internal/dto"
	"VitaTaskGo/internal/pkg/auth"
	"VitaTaskGo/internal/pkg/db"
	"VitaTaskGo/internal/pkg/response"
	"VitaTaskGo/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ProjectApi struct {
}

func NewProjectApi() *ProjectApi {
	return &ProjectApi{}
}

// CreateProject 创建项目信息
func (receiver *ProjectApi) CreateProject(ctx *gin.Context) {
	var post dto.CreateProjectForm
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}

	ctx.JSON(http.StatusOK, response.Auto(service.NewProjectService(db.Db, ctx).CreateProject(post.Name, post.Leader)))
}

// EditProject 编辑项目信息
func (receiver *ProjectApi) EditProject(ctx *gin.Context) {
	var post dto.EditProjectForm
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}
	ctx.JSON(http.StatusOK, response.Auto(service.NewProjectService(db.Db, ctx).EditProject(post.ID, post.Name, post.Leader)))
}

// ProjectList 项目列表
func (receiver *ProjectApi) ProjectList(ctx *gin.Context) {
	var post dto.ProjectListQuery
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}
	ctx.JSON(http.StatusOK, response.Auto(service.NewProjectService(db.Db, ctx).GetProjectList(post)))
}

// SimpleList 简单项目列表
func (receiver *ProjectApi) SimpleList(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, response.SuccessData(service.NewProjectService(db.Db, ctx).GetSimpleList()))
}

// ProjectTrash 项目回收站
func (*ProjectApi) ProjectTrash(ctx *gin.Context) {
	var post dto.ProjectListQuery
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}
	post.Deleted = true // 设置为搜索删除项目
	ctx.JSON(http.StatusOK, response.Auto(service.NewProjectService(db.Db, ctx).GetProjectList(post)))
}

// ProjectDelete 删除项目
func (receiver *ProjectApi) ProjectDelete(ctx *gin.Context) {
	var post dto.ProjectSingleId
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}
	if err := service.NewProjectService(db.Db, ctx).ProjectDelete(post.ID); err == nil {
		ctx.JSON(http.StatusOK, response.Success())
	} else {
		ctx.JSON(http.StatusOK, response.Error(err))
	}
}

// ProjectArchive 项目归档
func (receiver *ProjectApi) ProjectArchive(ctx *gin.Context) {
	var post dto.ProjectSingleId
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}

	ctx.JSON(http.StatusOK, response.Auto(nil, service.NewProjectService(db.Db, ctx).ProjectArchive(post.ID)))
}

// UnArchive 项目归档
func (receiver *ProjectApi) UnArchive(ctx *gin.Context) {
	var post dto.ProjectSingleId
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}

	ctx.JSON(http.StatusOK, response.Auto(nil, service.NewProjectService(db.Db, ctx).UnArchive(post.ID)))
}

// Star 收藏项目
func (receiver *ProjectApi) Star(ctx *gin.Context) {
	var post dto.ProjectSingleId
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}
	// 获取当前用户
	currUser, err := auth.CurrUser(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, response.Error(err))
		return
	}
	err = service.NewProjectMemberService(db.Db, ctx).ProjectStar(post.ID, currUser.ID)
	if err != nil {
		ctx.JSON(http.StatusOK, response.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, response.Success())
}

// UnStart 取消收藏项目
func (receiver *ProjectApi) UnStart(ctx *gin.Context) {
	var post dto.ProjectSingleId
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}
	// 获取当前用户
	currUser, err := auth.CurrUser(ctx)
	if err != nil {
		ctx.JSON(http.StatusOK, response.Error(err))
		return
	}
	err = service.NewProjectMemberService(db.Db, ctx).ProjectUnStar(post.ID, currUser.ID)
	if err != nil {
		ctx.JSON(http.StatusOK, response.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, response.Success())
}

// Transfer 移交项目
func (receiver *ProjectApi) Transfer(ctx *gin.Context) {
	var post dto.ProjectTransferForm
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}
	ctx.JSON(
		http.StatusOK,
		response.Auto(nil, service.NewProjectService(db.Db, ctx).Transfer(post.Project, post.Recipient)),
	)
}

// Detail 获取项目详情
func (receiver *ProjectApi) Detail(ctx *gin.Context) {
	var post dto.ProjectSingleId
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}

	ctx.JSON(
		http.StatusOK,
		response.Auto(service.NewProjectService(db.Db, ctx).GetOneProject(post.ID)),
	)
}
