package controller

import (
	"VitaTaskGo/app/logic"
	"VitaTaskGo/app/response"
	"VitaTaskGo/app/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

type TaskGroupController struct {
}

func NewTaskGroupController() *TaskGroupController {
	return &TaskGroupController{}
}

func (receiver TaskGroupController) Add(ctx *gin.Context) {
	var post types.TaskGroupForm
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}

	ctx.JSON(
		http.StatusOK,
		response.Auto(logic.NewTaskGroupLogic(ctx).Add(post)),
	)
}

func (receiver TaskGroupController) Update(ctx *gin.Context) {
	var post types.TaskGroupForm
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}

	ctx.JSON(
		http.StatusOK,
		response.Auto(logic.NewTaskGroupLogic(ctx).Update(post)),
	)
}

func (receiver TaskGroupController) Delete(ctx *gin.Context) {
	var post types.SingleUintRequired
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}

	ctx.JSON(
		http.StatusOK,
		response.Auto(nil, logic.NewTaskGroupLogic(ctx).Delete(post.ID)),
	)
}

func (receiver TaskGroupController) List(ctx *gin.Context) {
	var post types.TaskGroupQuery
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}

	ctx.JSON(
		http.StatusOK,
		response.SuccessData(logic.NewTaskGroupLogic(ctx).List(post)),
	)
}

func (receiver TaskGroupController) Detail(ctx *gin.Context) {
	var post types.SingleUintRequired
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}

	ctx.JSON(
		http.StatusOK,
		response.Auto(logic.NewTaskGroupLogic(ctx).Detail(post.ID)),
	)
}

func (receiver TaskGroupController) SimpleList(ctx *gin.Context) {
	var post types.UintId
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}

	if post.ID > 0 {
		ctx.JSON(
			http.StatusOK,
			response.SuccessData(logic.NewTaskGroupLogic(ctx).SimpleList(post.ID)),
		)
	} else {
		// 如果没有提供项目id，查询全部项目并且按项目分组
		ctx.JSON(
			http.StatusOK,
			response.SuccessData(logic.NewTaskGroupLogic(ctx).SimpleGroupList()),
		)
	}
}
