package handle

import (
	"VitaTaskGo/internal/api/model/dto"
	"VitaTaskGo/internal/api/service"
	"VitaTaskGo/pkg/db"
	"VitaTaskGo/pkg/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

type TaskGroupApi struct {
}

func NewTaskGroupApi() *TaskGroupApi {
	return &TaskGroupApi{}
}

func (receiver TaskGroupApi) Add(ctx *gin.Context) {
	var post dto.TaskGroupForm
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}

	ctx.JSON(
		http.StatusOK,
		response.Auto(service.NewTaskGroupService(db.Db, ctx).Add(post)),
	)
}

func (receiver TaskGroupApi) Update(ctx *gin.Context) {
	var post dto.TaskGroupForm
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}

	ctx.JSON(
		http.StatusOK,
		response.Auto(service.NewTaskGroupService(db.Db, ctx).Update(post)),
	)
}

func (receiver TaskGroupApi) Delete(ctx *gin.Context) {
	var post dto.SingleUintRequired
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}

	ctx.JSON(
		http.StatusOK,
		response.Auto(nil, service.NewTaskGroupService(db.Db, ctx).Delete(post.ID)),
	)
}

func (receiver TaskGroupApi) List(ctx *gin.Context) {
	var post dto.TaskGroupQuery
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}

	ctx.JSON(
		http.StatusOK,
		response.SuccessData(service.NewTaskGroupService(db.Db, ctx).List(post)),
	)
}

func (receiver TaskGroupApi) Detail(ctx *gin.Context) {
	var post dto.SingleUintRequired
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}

	ctx.JSON(
		http.StatusOK,
		response.Auto(service.NewTaskGroupService(db.Db, ctx).Detail(post.ID)),
	)
}

func (receiver TaskGroupApi) SimpleList(ctx *gin.Context) {
	var post dto.UintId
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}

	if post.ID > 0 {
		ctx.JSON(
			http.StatusOK,
			response.SuccessData(service.NewTaskGroupService(db.Db, ctx).SimpleList(post.ID)),
		)
	} else {
		// 如果没有提供项目id，查询全部项目并且按项目分组
		ctx.JSON(
			http.StatusOK,
			response.SuccessData(service.NewTaskGroupService(db.Db, ctx).SimpleGroupList()),
		)
	}
}
