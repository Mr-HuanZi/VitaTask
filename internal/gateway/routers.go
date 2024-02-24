package gateway

import (
	"VitaTaskGo/internal/gateway/handle"
	"github.com/gin-gonic/gin"
)

func Routers(r *gin.Engine) {
	group := r.Group("gateway")

	{
		registerApi := handle.NewRegisterApi()
		group.POST("bind", registerApi.BindUser)
	}
}

func WebSocketRouters(r *gin.Engine) {
	r.GET("chat", handle.ClientHandle)
}
