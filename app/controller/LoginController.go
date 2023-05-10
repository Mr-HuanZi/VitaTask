package controller

import (
	"VitaTaskGo/app/logic"
	"VitaTaskGo/app/modules/ws"
	"VitaTaskGo/app/response"
	"VitaTaskGo/app/types"
	"github.com/gin-gonic/gin"
	"net/http"
)

type LoginController struct {
}

func NewLoginController() *LoginController {
	return &LoginController{}
}

// Login 登录接口
// Api POST /login
func (*LoginController) Login(ctx *gin.Context) {
	var (
		post types.LoginForm
	)
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}
	token, user, err := logic.NewLoginLogic(ctx).UserLogin(post.Username, post.Password)
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
func (*LoginController) Register(ctx *gin.Context) {
	var (
		post types.UserRegisterForm
	)
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}
	err := logic.NewLoginLogic(ctx).UserRegister(post)
	ctx.JSON(http.StatusOK, response.Auto(nil, err))
}
