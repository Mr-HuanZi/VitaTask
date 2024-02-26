package gateway

import (
	"VitaTaskGo/internal/gateway/handle"
	"github.com/gin-gonic/gin"
)

func Routers(r *gin.Engine) {
	group := r.Group("gateway")

	{
		chatApi := handle.NewChatApi()
		group.POST("send/user", chatApi.SendToUser)
		group.POST("send/users", chatApi.SendToUsers)
	}
}

func WebSocketRouters(r *gin.Engine) {
	r.GET("chat", handle.ClientHandle)
}
