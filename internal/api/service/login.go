package service

import (
	"VitaTaskGo/internal/api/data"
	"VitaTaskGo/internal/api/model/dto"
	"VitaTaskGo/internal/pkg"
	"VitaTaskGo/internal/pkg/auth"
	"VitaTaskGo/internal/repo"
	"VitaTaskGo/pkg/exception"
	"VitaTaskGo/pkg/response"
	"errors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"time"
)

type LoginService struct {
	Db   *gorm.DB
	ctx  *gin.Context
	repo repo.UserRepo
}

func NewLoginService(tx *gorm.DB, ctx *gin.Context) *LoginService {
	return &LoginService{
		Db:   tx,  // 赋予ORM实例
		ctx:  ctx, // 传递上下文
		repo: data.NewUserRepo(tx, ctx),
	}
}

func (s LoginService) UserLogin(username, password string) (string, *repo.User, error) {
	// 查询用户
	user, err := s.repo.QueryUsernameAndPass(username, pkg.Encryption(password)) // 用户密码加密
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// 用户名或密码不正确
			return "", nil, exception.NewException(response.LoginPassError)
		}
		return "", nil, exception.ErrorHandle(err, response.DbQueryError)
	}

	token, err := auth.GenerateToken(user.ID, user.UserLogin)
	if err != nil {
		return "", nil, exception.NewException(response.LoginSingGenerateFail)
	}
	// 生成Token
	return token, user, nil
}

func (s LoginService) UserRegister(post dto.UserRegisterForm) error {
	// 查询用户名
	_, err := s.repo.QueryUsername(post.Username)
	if err == nil {
		// 没有错误表示查询到了记录，说明输入的用户名已被占用
		// 用户名已存在
		return exception.NewException(response.RegUsernameExists)
	} else {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			// 查询错误
			return exception.ErrorHandle(err, response.DbQueryError)
		}
	}

	// 组建用户数据
	newUser := &repo.User{
		UserStatus:   1,
		UserLogin:    post.Username,
		UserPass:     pkg.Encryption(post.Password), // 密码加密
		UserNickname: post.UserNickname,
		UserEmail:    post.UserEmail,
		Mobile:       post.Mobile,
		LockTime:     0,
		ErrorSum:     0,
		First:        1,                 // 是否首次登录
		LastEditPass: time.Now().Unix(), // 最后一次修改密码时间，记录为当前
	}
	// 插入数据
	createErr := s.repo.CreateUser(newUser)
	if createErr != nil {
		// 写入数据失败
		return exception.ErrorHandle(err, response.RegFail)
	}
	return nil
}
