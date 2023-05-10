package controller

import (
	"VitaTaskGo/app/logic"
	"VitaTaskGo/app/response"
	"VitaTaskGo/app/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

type DialogController struct {
}

func NewDialogController() *DialogController {
	return &DialogController{}
}

// List 对话列表
func (receiver DialogController) List(ctx *gin.Context) {
	//
}

// MsgList 消息列表
func (receiver DialogController) MsgList(ctx *gin.Context) {
	var dto types.DialogIdDto
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}

	ctx.JSON(http.StatusOK, response.Auto(logic.NewDialogLogic(ctx).MsgList(dto.DialogId)))
}

// SendText 发送文本消息
func (receiver DialogController) SendText(ctx *gin.Context) {
	var dto types.DialogSendTextDto
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}

	ctx.JSON(http.StatusOK, response.Auto(logic.NewDialogLogic(ctx).SendText(dto)))
}

func (receiver DialogController) Create(ctx *gin.Context) {
	var dto types.DialogCreateDto
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}

	ctx.JSON(http.StatusOK, response.Auto(logic.NewDialogLogic(ctx).Create(dto.Name, dto.Type, dto.Members)))
}
