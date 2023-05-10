package logic

import (
	"VitaTaskGo/app/exception"
	"VitaTaskGo/app/extend"
	"VitaTaskGo/app/extend/jwt"
	"VitaTaskGo/app/model"
	"VitaTaskGo/app/response"
	"VitaTaskGo/app/types"
	"VitaTaskGo/library/db"
	"github.com/gin-gonic/gin"
	"time"
)

type LoginLogic struct {
	ctx *gin.Context
}

func NewLoginLogic(ctx *gin.Context) *LoginLogic {
	return &LoginLogic{
		ctx: ctx, // 传递上下文
	}
}

func (s LoginLogic) UserLogin(username, password string) (string, *model.User, error) {
	// 用户密码加密
	pwd := extend.Encryption(password)
	user := &model.User{}
	// 查询用户
	db.Db.Where("user_login = ?", username).Where("user_pass = ?", pwd).First(user)
	if user.ID <= 0 {
		// 用户名或密码不正确
		return "", nil, exception.NewException(response.LoginPassError)
	}

	token, err := jwt.GenerateToken(user.ID, user.UserLogin)
	if err != nil {
		return "", nil, exception.NewException(response.LoginSingGenerateFail)
	}
	// 生成Token
	return token, user, nil
}

func (s LoginLogic) UserRegister(post types.UserRegisterForm) error {
	var count int64
	db.Db.Model(&model.User{}).Where("user_login = ?", post.Username).Count(&count)
	if count > 0 {
		// 用户名已存在
		return exception.NewException(response.RegUsernameExists)
	}

	// 组建用户数据
	newUser := &model.User{
		UserStatus:   1,
		UserLogin:    post.Username,
		UserPass:     extend.Encryption(post.Password), // 密码加密
		UserNickname: post.UserNickname,
		UserEmail:    post.UserEmail,
		Mobile:       post.Mobile,
		LockTime:     0,
		ErrorSum:     0,
		First:        1,                 // 是否首次登录
		LastEditPass: time.Now().Unix(), // 最后一次修改密码时间，记录为当前
	}
	// 插入数据
	result := db.Db.Create(newUser)
	if result.Error != nil {
		// 写入数据失败
		return exception.NewException(response.RegFail)
	}
	return nil
}
