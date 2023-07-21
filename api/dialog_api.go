package api

import (
	"VitaTaskGo/internal/dto"
	"VitaTaskGo/internal/pkg/db"
	"VitaTaskGo/internal/pkg/response"
	"VitaTaskGo/internal/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

type DialogApi struct {
}

func NewDialogApi() *DialogApi {
	return &DialogApi{}
}

// MsgList 消息列表
func (receiver DialogApi) MsgList(ctx *gin.Context) {
	var dialogIdDto dto.DialogIdDto
	if err := ctx.ShouldBindJSON(&dialogIdDto); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}

	ctx.JSON(http.StatusOK, response.Auto(service.NewDialogService(db.Db, ctx).MsgList(dialogIdDto.DialogId)))
}

// SendText 发送文本消息
func (receiver DialogApi) SendText(ctx *gin.Context) {
	var dialogSendTextDto dto.DialogSendTextDto
	if err := ctx.ShouldBindJSON(&dialogSendTextDto); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}

	ctx.JSON(http.StatusOK, response.Auto(service.NewDialogService(db.Db, ctx).SendText(dialogSendTextDto)))
}

func (receiver DialogApi) Create(ctx *gin.Context) {
	var dialogCreateDto dto.DialogCreateDto
	if err := ctx.ShouldBindJSON(&dialogCreateDto); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}

	ctx.JSON(http.StatusOK, response.Auto(service.NewDialogService(db.Db, ctx).Create(dialogCreateDto.Name, dialogCreateDto.Type, dialogCreateDto.Members)))
}
