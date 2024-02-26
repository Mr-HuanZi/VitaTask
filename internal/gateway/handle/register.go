package handle

import (
	"VitaTaskGo/internal/pkg/gateway"
	"VitaTaskGo/pkg/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ClientHandle(ctx *gin.Context) {
	// 创建客户端实例
	client := gateway.NewChatClient()
	// 创建连接
	err := client.Conn(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, response.Exception(response.SystemFail))
		return
	}
}
