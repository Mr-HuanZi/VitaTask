package jwt

import (
	"VitaTaskGo/app/exception"
	"VitaTaskGo/app/response"
	"VitaTaskGo/library/config"
	jwtGo "github.com/dgrijalva/jwt-go"
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
