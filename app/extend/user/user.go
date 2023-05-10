package user

import (
	"VitaTaskGo/app/constant"
	"VitaTaskGo/app/exception"
	"VitaTaskGo/app/extend/jwt"
	"VitaTaskGo/app/model"
	"VitaTaskGo/app/response"
	"VitaTaskGo/library/db"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"strings"
)

// ParseAuthorization 解析Authorization
func ParseAuthorization(authorization string) (*jwt.MyCustomClaims, error) {
	if authorization == "" {
		return nil, exception.NewException(response.SignatureMissing)
	}
	// 检查字符串开头是否包含 “Bearer ”
	if strings.HasPrefix(authorization, "Bearer") {
		authorization = strings.TrimSpace(strings.TrimPrefix(authorization, "Bearer"))
	}
	return jwt.ParseToken(authorization)
}

// CurrUser 获取当前登录用户
// 如果用户被禁用会返回错误
func CurrUser(ctx *gin.Context) (*model.User, error) {
	var user *model.User
	currUid, ok := ctx.Get(constant.CurrUidKey)
	if !ok {
		return nil, exception.NewException(response.NotLoggedIn)
	}
	err := db.Db.First(&user, currUid).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// 用户不存在
			return nil, exception.NewException(response.UserNotFound)
		} else {
			// 其它错误
			return nil, exception.ErrorHandle(err, response.DbQueryError)
		}
	}
	// 检查用户是否被禁用
	if user.UserStatus != 1 {
		return nil, exception.NewException(response.UserDisabled)
	}
	return user, err
}

// IsSuper 是否超级账户
func IsSuper(user *model.User) bool {
	return user.Super == 1
}

// Exist 用户是否存在
func Exist(uid uint64) bool {
	var count int64
	if err := db.Db.Model(&model.User{}).Where("id = ?", uid).Count(&count).Error; err != nil {
		logrus.Errorln(err)
		return false
	}
	return count > 0
}
