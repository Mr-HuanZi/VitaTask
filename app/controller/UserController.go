package controller

import (
	"VitaTaskGo/app/logic"
	"VitaTaskGo/app/response"
	"VitaTaskGo/app/types"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

type UserController struct {
}

func NewUserController() *UserController {
	return &UserController{}
}

// CurrUser 获取当前登录用户
func (receiver UserController) CurrUser(ctx *gin.Context) {
	currUser, err := logic.NewUserLogic(ctx).CurrUser()
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, response.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, response.SuccessData(currUser))
}

func (receiver UserController) StoreSelf(ctx *gin.Context) {
	var post types.UserInfoDto
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}
	// 验证日期格式是否正确
	if len(strings.TrimSpace(post.Birthday)) > 0 {
		_, err := time.ParseInLocation(time.DateOnly, strings.TrimSpace(post.Birthday), time.Local)
		if err != nil {
			ctx.JSON(http.StatusOK, response.Exception(response.TimeParseFail))
			return
		}
	}

	ctx.JSON(http.StatusOK, response.Auto(logic.NewUserLogic(ctx).StoreSelf(post)))
}

func (receiver UserController) ChangeAvatar(ctx *gin.Context) {
	var post types.FileDto
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}

	ctx.JSON(http.StatusOK, response.Auto(nil, logic.NewUserLogic(ctx).ChangeAvatar(post)))
}

func (receiver UserController) ChangePassword(ctx *gin.Context) {
	var post types.ChangePasswordDto
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}

	ctx.JSON(http.StatusOK, response.Auto(nil, logic.NewUserLogic(ctx).ChangePassword(post)))
}

func (receiver UserController) ChangeMobile(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, response.Auto(nil, logic.NewUserLogic(ctx).ChangeMobile(ctx.Query("mobile"))))
}

func (receiver UserController) ChangeEmail(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, response.Auto(nil, logic.NewUserLogic(ctx).ChangeEmail(ctx.Query("email"))))
}
