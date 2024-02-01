package handle

import (
	"VitaTaskGo/internal/api/model/dto"
	"VitaTaskGo/internal/api/service"
	"VitaTaskGo/internal/pkg/ws"
	"VitaTaskGo/pkg/db"
	"VitaTaskGo/pkg/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

type LoginApi struct {
}

func NewLoginApi() *LoginApi {
	return &LoginApi{}
}

// Login 登录接口
// Api POST /login
func (*LoginApi) Login(ctx *gin.Context) {
	var (
		post dto.LoginForm
	)
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}
	token, user, err := service.NewLoginService(db.Db, ctx).UserLogin(post.Username, post.Password)
	if err != nil {
		ctx.JSON(http.StatusOK, response.Error(err))
		return
	}

	ctx.JSON(http.StatusOK, response.SuccessData(map[string]interface{}{
		"token":         token,
		"id":            user.ID,
		"user_nickname": user.UserNickname,
		"user_login":    user.UserLogin,
		// 生成websocket需要的Token，一次性的，每次登录后重新生成
		"ws_token": ws.GenerateToken([]string{user.UserLogin, user.UserLogin}),
	}))
}

// Register 注册接口
// Api POST /register
func (*LoginApi) Register(ctx *gin.Context) {
	var (
		post dto.UserRegisterForm
	)
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}
	err := service.NewLoginService(db.Db, ctx).UserRegister(post)
	ctx.JSON(http.StatusOK, response.Auto(nil, err))
}
