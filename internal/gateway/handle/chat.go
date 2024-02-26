package handle

import (
	"VitaTaskGo/internal/gateway/model/dto"
	"VitaTaskGo/internal/gateway/services"
	"VitaTaskGo/pkg/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

type ChatApi struct {
}

func NewChatApi() *ChatApi {
	return &ChatApi{}
}

func (receiver *ChatApi) SendToUser(ctx *gin.Context) {
	var post dto.ChatSendUserForm
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}
	ctx.JSON(http.StatusOK, response.Auto(nil, services.NewChatService(ctx).SendToUser(post.Userid, post.Msg)))
}

func (receiver *ChatApi) SendToUsers(ctx *gin.Context) {
	var post dto.ChatSendUsersForm
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}
	ctx.JSON(http.StatusOK, response.Auto(nil, services.NewChatService(ctx).SendToUsers(post.Users, post.Msg)))
}
