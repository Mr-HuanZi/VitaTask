package api

import (
	"VitaTaskGo/internal/dto"
	"VitaTaskGo/internal/pkg"
	"VitaTaskGo/internal/pkg/db"
	"VitaTaskGo/internal/pkg/response"
	"VitaTaskGo/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type TaskApi struct {
}

func NewTaskApi() *TaskApi {
	return &TaskApi{}
}

func (receiver TaskApi) Lists(ctx *gin.Context) {
	var query dto.TaskListQuery
	if err := ctx.ShouldBindJSON(&query); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}

	ctx.JSON(
		http.StatusOK,
		response.Auto(service.NewTaskService(db.Db, ctx).Lists(query)),
	)
}

func (receiver TaskApi) Create(ctx *gin.Context) {
	var post dto.TaskCreateForm
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}

	ctx.JSON(
		http.StatusOK,
		response.Auto(service.NewTaskService(db.Db, ctx).Create(post)),
	)
}

func (receiver TaskApi) Detail(ctx *gin.Context) {
	var post dto.SingleUintRequired
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}

	ctx.JSON(
		http.StatusOK,
		response.Auto(service.NewTaskService(db.Db, ctx).Detail(post.ID)),
	)
}

func (receiver TaskApi) Roles(ctx *gin.Context) {
	ctx.JSON(
		http.StatusOK,
		response.SuccessData(service.NewTaskService(db.Db, ctx).Roles()),
	)
}

func (receiver TaskApi) Status(ctx *gin.Context) {
	ctx.JSON(
		http.StatusOK,
		response.SuccessData(service.NewTaskService(db.Db, ctx).Status()),
	)
}

func (receiver TaskApi) ChangeStatus(ctx *gin.Context) {
	var post dto.TaskChangeStatus
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}

	ctx.JSON(
		http.StatusOK,
		response.Auto(nil, service.NewTaskService(db.Db, ctx).ChangeStatus(post.ID, post.Status)),
	)
}

func (receiver TaskApi) Update(ctx *gin.Context) {
	var post dto.TaskCreateForm
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}
	// 从Query中获取任务id
	id := pkg.ParseStringToUi64(ctx.Query("id"))

	ctx.JSON(
		http.StatusOK,
		response.Auto(service.NewTaskService(db.Db, ctx).Update(uint(id), post)),
	)
}

func (receiver TaskApi) Delete(ctx *gin.Context) {
	var post dto.SingleUintRequired
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}

	ctx.JSON(
		http.StatusOK,
		response.Auto(nil, service.NewTaskService(db.Db, ctx).Delete(post.ID)),
	)
}

func (receiver TaskApi) Statistics(ctx *gin.Context) {
	var post dto.SingleUintRequired // 项目id
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}

	ctx.JSON(
		http.StatusOK,
		response.SuccessData(service.NewTaskService(db.Db, ctx).Statistics(post.ID)),
	)
}

func (receiver TaskApi) DailySituation(ctx *gin.Context) {
	var post dto.DailySituationQuery
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}

	ctx.JSON(
		http.StatusOK,
		response.Auto(service.NewTaskService(db.Db, ctx).DailySituation(post)),
	)
}
