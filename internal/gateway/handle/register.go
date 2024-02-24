package handle

import "github.com/gin-gonic/gin"

type RegisterApi struct {
}

func NewRegisterApi() *RegisterApi {
	return &RegisterApi{}
}

func (receiver *RegisterApi) BindUser(ctx *gin.Context) {

}

func (receiver *RegisterApi) SendToUser(ctx *gin.Context) {

}

func (receiver *RegisterApi) SendToUsers(ctx *gin.Context) {

}
