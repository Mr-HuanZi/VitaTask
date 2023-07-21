package auth

import (
	"VitaTaskGo/internal/biz"
	"VitaTaskGo/internal/constant"
	"VitaTaskGo/internal/data"
	"VitaTaskGo/internal/pkg"
	"VitaTaskGo/internal/pkg/config"
	"VitaTaskGo/internal/pkg/db"
	"VitaTaskGo/internal/pkg/exception"
	"VitaTaskGo/internal/pkg/response"
	jwtGo "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strings"
	"time"
)

type MyCustomClaims struct {
	UserId   uint64
	Username string
	jwtGo.StandardClaims
}

// GenerateToken 生成Token
func GenerateToken(userId uint64, username string) (string, error) {
	expireSeconds := config.Instances.Jwt.ExpireSeconds
	if expireSeconds <= 0 {
		// 默认10分钟过期
		expireSeconds = 600
	}
	expiresAt := time.Now().Add(time.Second * time.Duration(expireSeconds)).Unix()

	newClaims := MyCustomClaims{
		UserId:   userId,
		Username: username,
		StandardClaims: jwtGo.StandardClaims{
			// 过期时间
			ExpiresAt: expiresAt,
			// 签发时间
			IssuedAt: time.Now().Unix(),
			// 签发人
			Issuer: config.Instances.Jwt.Issuer,
		},
	}
	token := jwtGo.NewWithClaims(jwtGo.SigningMethodHS256, newClaims)

	tokenString, err := token.SignedString([]byte(config.Instances.Jwt.Key))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// ParseToken 解析Token
func ParseToken(tokenString string) (*MyCustomClaims, error) {
	if tokenString == "" {
		return nil, exception.NewException(response.SignatureMissing)
	}

	claims := &MyCustomClaims{} // 将Claims解析到这个结构体
	_, err := jwtGo.ParseWithClaims(tokenString, claims, func(token *jwtGo.Token) (interface{}, error) {
		return []byte(config.Instances.Jwt.Key), nil
	})

	if err != nil {
		return nil, err
	}
	return claims, nil
}

// ParseAuthorization 解析Authorization
func ParseAuthorization(authorization string) (*MyCustomClaims, error) {
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
func CurrUser(ctx *gin.Context) (*biz.User, error) {
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
func IsSuper(user *biz.User) bool {
	return user.Super == 1
}
