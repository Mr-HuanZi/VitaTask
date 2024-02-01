package handle

import (
	"VitaTaskGo/internal/api/model/dto"
	"VitaTaskGo/internal/api/service"
	"VitaTaskGo/pkg/db"
	"VitaTaskGo/pkg/response"
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
