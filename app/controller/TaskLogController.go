package controller

import (
	"VitaTaskGo/app/logic"
	"VitaTaskGo/app/response"
	"VitaTaskGo/app/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

type TaskLogController struct {
}

func NewTaskLogController() *TaskLogController {
	return &TaskLogController{}
}

func (receiver TaskLogController) List(ctx *gin.Context) {
	var query types.TaskLogQuery
	if err := ctx.ShouldBindJSON(&query); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}

	ctx.JSON(
		http.StatusOK,
		response.SuccessData(logic.NewTaskLogLogic(ctx).List(query)),
	)
}

func (receiver TaskLogController) Operators(ctx *gin.Context) {
	ctx.JSON(
		http.StatusOK,
		response.SuccessData(logic.NewTaskLogLogic(ctx).Operators()),
	)
}
