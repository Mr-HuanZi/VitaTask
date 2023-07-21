package api

import (
	"VitaTaskGo/internal/dto"
	"VitaTaskGo/internal/pkg/db"
	"VitaTaskGo/internal/pkg/response"
	"VitaTaskGo/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type TaskLogApi struct {
}

func NewTaskLogApi() *TaskLogApi {
	return &TaskLogApi{}
}

func (receiver TaskLogApi) List(ctx *gin.Context) {
	var query dto.TaskLogQuery
	if err := ctx.ShouldBindJSON(&query); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}

	ctx.JSON(
		http.StatusOK,
		response.SuccessData(service.NewTaskLogService(db.Db, ctx).List(query)),
	)
}

func (receiver TaskLogApi) Operators(ctx *gin.Context) {
	ctx.JSON(
		http.StatusOK,
		response.SuccessData(service.NewTaskLogService(db.Db, ctx).Operators()),
	)
}
