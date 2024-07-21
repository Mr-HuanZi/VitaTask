package auth

import (
	"VitaTaskGo/internal/api/data"
	"VitaTaskGo/internal/pkg"
	"VitaTaskGo/internal/pkg/constant"
	"VitaTaskGo/internal/repo"
	"VitaTaskGo/pkg/config"
	"VitaTaskGo/pkg/db"
	"VitaTaskGo/pkg/exception"
	"VitaTaskGo/pkg/response"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
	"strings"
	"time"
)

type UserJwtClaims struct {
	UserId   uint64
	Username string
	jwt.RegisteredClaims
}

// GenerateToken 生成Token
func GenerateToken(userId uint64, username string) (string, error) {
	expireSeconds := config.Get().Jwt.ExpireSeconds
	if expireSeconds <= 0 {
		// 默认10分钟过期
		expireSeconds = 600
	}
	expiresAt := time.Now().Add(time.Second * time.Duration(expireSeconds))

	newClaims := UserJwtClaims{
		userId,
		username,
		jwt.RegisteredClaims{
			Issuer:    config.Get().Jwt.Issuer,        // 签发人
			ExpiresAt: jwt.NewNumericDate(expiresAt),  // 过期时间
			IssuedAt:  jwt.NewNumericDate(time.Now()), // 签发时间
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, newClaims)

	tokenString, err := token.SignedString([]byte(config.Get().Jwt.Key))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ParseToken 解析Token
func ParseToken(tokenString string) (*UserJwtClaims, error) {
	if tokenString == "" {
		return nil, exception.NewException(response.SignatureMissing)
	}

	claims := &UserJwtClaims{} // 将Claims解析到这个结构体
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Get().Jwt.Key), nil
	})

	if err != nil {
		return nil, err
	}
	return claims, nil
}

// ParseAuthorization 解析Authorization
func ParseAuthorization(authorization string) (*UserJwtClaims, error) {
	if authorization == "" {
		return nil, exception.NewException(response.SignatureMissing)
	}
	// 检查字符串开头是否包含 “Bearer ”
	if strings.HasPrefix(authorization, "Bearer") {
		authorization = strings.TrimSpace(strings.TrimPrefix(authorization, "Bearer"))
	}
	return ParseToken(authorization)
}

// CurrUser 获取当前登录用户
// 如果用户被禁用会返回错误
func CurrUser(ctx *gin.Context) (*repo.User, error) {
	var uid uint64

	currUid, ok := ctx.Get(constant.CurrUidKey)
	if !ok {
		return nil, exception.NewException(response.NotLoggedIn)
	}

	// 处理类型转换
	switch currUid.(type) {
	case uint64:
		uid = currUid.(uint64)
		break
	case int64:
		uid = uint64(currUid.(int64))
	case string:
		uid = pkg.ParseStringToUi64(currUid.(string))
		break
	default:
		return nil, exception.NewException(response.NotLoggedIn)
	}

	user, err := data.NewUserRepo(db.Db, ctx).GetUser(uid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
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
func IsSuper(user *repo.User) bool {
	return user.Super == 1
}
