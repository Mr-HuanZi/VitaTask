package handle

import (
	"VitaTaskGo/internal/api/model/dto"
	"VitaTaskGo/internal/api/service"
	"VitaTaskGo/pkg/db"
	"VitaTaskGo/pkg/response"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
	"time"
)

type UserApi struct {
}

func NewUserApi() *UserApi {
	return &UserApi{}
}

// CurrUser 获取当前登录用户
func (receiver UserApi) CurrUser(ctx *gin.Context) {
	currUser, err := service.NewUserService(db.Db, ctx).CurrUser()
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, response.Error(err))
		return
	}
	ctx.JSON(http.StatusOK, response.SuccessData(currUser))
}

func (receiver UserApi) StoreSelf(ctx *gin.Context) {
	var post dto.UserInfoDto
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

	ctx.JSON(http.StatusOK, response.Auto(service.NewUserService(db.Db, ctx).StoreSelf(post)))
}

func (receiver UserApi) ChangeAvatar(ctx *gin.Context) {
	var post dto.FileDto
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}

	ctx.JSON(http.StatusOK, response.Auto(nil, service.NewUserService(db.Db, ctx).ChangeAvatar(post)))
}

func (receiver UserApi) ChangePassword(ctx *gin.Context) {
	var post dto.ChangePasswordDto
	if err := ctx.ShouldBindJSON(&post); err != nil {
		ctx.JSON(http.StatusOK, response.HandleFormVerificationFailed(err))
		return
	}

	ctx.JSON(http.StatusOK, response.Auto(nil, service.NewUserService(db.Db, ctx).ChangePassword(post)))
}

func (receiver UserApi) ChangeMobile(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, response.Auto(nil, service.NewUserService(db.Db, ctx).ChangeMobile(ctx.Query("mobile"))))
}

func (receiver UserApi) ChangeEmail(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, response.Auto(nil, service.NewUserService(db.Db, ctx).ChangeEmail(ctx.Query("email"))))
}
