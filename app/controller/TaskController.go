package controller

import (
	"VitaTaskGo/app/extend"
	"VitaTaskGo/app/logic"
	"VitaTaskGo/app/response"
	"VitaTaskGo/app/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

type TaskController struct {
}

func NewTaskController() *TaskController {
	return &TaskController{}
}

func (receiver TaskController) Lists(ctx *gin.Context) {
	var query types.TaskListQuery
	if err := ctx.ShouldBindJSON(&query); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}

	ctx.JSON(
		http.StatusOK,
		response.Auto(logic.NewTaskLogic(ctx).Lists(query)),
	)
}

func (receiver TaskController) Create(ctx *gin.Context) {
	var post types.TaskCreateForm
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}

	ctx.JSON(
		http.StatusOK,
		response.Auto(logic.NewTaskLogic(ctx).Create(post)),
	)
}

func (receiver TaskController) Detail(ctx *gin.Context) {
	var post types.SingleUintRequired
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}

	ctx.JSON(
		http.StatusOK,
		response.Auto(logic.NewTaskLogic(ctx).Detail(post.ID)),
	)
}

func (receiver TaskController) Roles(ctx *gin.Context) {
	ctx.JSON(
		http.StatusOK,
		response.SuccessData(logic.NewTaskLogic(ctx).Roles()),
	)
}

func (receiver TaskController) Status(ctx *gin.Context) {
	ctx.JSON(
		http.StatusOK,
		response.SuccessData(logic.NewTaskLogic(ctx).Status()),
	)
}

func (receiver TaskController) ChangeStatus(ctx *gin.Context) {
	var post types.TaskChangeStatus
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}

	ctx.JSON(
		http.StatusOK,
		response.Auto(nil, logic.NewTaskLogic(ctx).ChangeStatus(post.ID, post.Status)),
	)
}

func (receiver TaskController) Update(ctx *gin.Context) {
	var post types.TaskCreateForm
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}
	// 从Query中获取任务id
	id := extend.ParseStringToUi64(ctx.Query("id"))

	ctx.JSON(
		http.StatusOK,
		response.Auto(logic.NewTaskLogic(ctx).Update(uint(id), post)),
	)
}

func (receiver TaskController) Delete(ctx *gin.Context) {
	var post types.SingleUintRequired
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}

	ctx.JSON(
		http.StatusOK,
		response.Auto(nil, logic.NewTaskLogic(ctx).Delete(post.ID)),
	)
}

func (receiver TaskController) Statistics(ctx *gin.Context) {
	var post types.SingleUintRequired // 项目id
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}

	ctx.JSON(
		http.StatusOK,
		response.SuccessData(logic.NewTaskLogic(ctx).Statistics(post.ID)),
	)
}

func (receiver TaskController) DailySituation(ctx *gin.Context) {
	var post types.DailySituationQuery
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}

	ctx.JSON(
		http.StatusOK,
		response.Auto(logic.NewTaskLogic(ctx).DailySituation(post)),
	)
}
