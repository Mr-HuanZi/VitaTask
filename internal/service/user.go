package service

import (
	"VitaTaskGo/internal/biz"
	"VitaTaskGo/internal/constant"
	"VitaTaskGo/internal/data"
	"VitaTaskGo/internal/dto"
	"VitaTaskGo/internal/pkg"
	"VitaTaskGo/internal/pkg/auth"
	"VitaTaskGo/internal/pkg/exception"
	"VitaTaskGo/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"github.com/gotidy/copy"
	"gorm.io/gorm"
	"strings"
	"time"
)

type UserService struct {
	Db   *gorm.DB
	ctx  *gin.Context
	repo biz.UserRepo
}

func NewUserService(tx *gorm.DB, ctx *gin.Context) *UserService {
	return &UserService{
		Db:   tx,  // 赋予ORM实例
		ctx:  ctx, // 传递上下文
		repo: data.NewUserRepo(tx, ctx),
	}
}

// CurrUser 返回当前登录用户
// 给API调用的
func (receiver UserService) CurrUser() (*biz.User, error) {
	return auth.CurrUser(receiver.ctx)
}

// UserExist 用户是否存在
func (receiver UserService) UserExist(uid uint64) bool {
	return receiver.repo.Exist(uid)
}

// StoreSelf 保存用户信息
func (receiver UserService) StoreSelf(data dto.UserInfoDto) (*biz.User, error) {
	// 获取登录用户uid
	uid, _ := receiver.ctx.Get(constant.CurrUidKey)
	if !receiver.UserExist(uid.(uint64)) {
		return nil, exception.NewException(response.UserNotFound)
	}
	// 拷贝变量
	userData := biz.User{}
	copiers := copy.New(func(c *copy.Options) {
		c.Skip = true
	})
	copiers.Copy(&userData, &data)

	err := receiver.repo.UpdatesUser(uid.(uint64), &userData)

	// 重新获取用户信息
	user, err := receiver.repo.GetUser(uid.(uint64))
	return user, err
}

// ChangeAvatar 变更头像
func (receiver UserService) ChangeAvatar(avatar dto.FileDto) error {
	currUser, err := auth.CurrUser(receiver.ctx)
	if err != nil {
		return err
	}
	// 是否提供url
	if len(strings.TrimSpace(avatar.Url)) <= 0 {
		return exception.NewException(response.AvatarNotUploaded)
	}

	currUser.Avatar = avatar.Url
	err = receiver.repo.SaveUser(currUser)
	return exception.ErrorHandle(err, response.DbExecuteError)
}

// ChangePassword 修改当前用户密码
func (receiver UserService) ChangePassword(data dto.ChangePasswordDto) error {
	currUser, err := auth.CurrUser(receiver.ctx)
	if err != nil {
		return err
	}

	// 旧密码匹配
	if pkg.Encryption(data.OldPassword) != currUser.UserPass {
		return exception.NewException(response.PassError)
	}

	currUser.UserPass = pkg.Encryption(data.Password)
	// 记录密码修改的时间
	currUser.LastEditPass = time.Now().UnixMilli()
	err = receiver.repo.SaveUser(currUser)
	return exception.ErrorHandle(err, response.DbExecuteError)
}

// ChangeMobile 变更手机号
func (receiver UserService) ChangeMobile(mobile string) error {
	currUser, err := auth.CurrUser(receiver.ctx)
	if err != nil {
		return err
	}
	if len(strings.TrimSpace(mobile)) <= 0 {
		return exception.NewException(response.NotInputtedMobile)
	}

	currUser.Mobile = mobile
	err = receiver.repo.SaveUser(currUser)
	return exception.ErrorHandle(err, response.DbExecuteError)
}

// ChangeEmail 变更电子邮箱地址
func (receiver UserService) ChangeEmail(email string) error {
	currUser, err := auth.CurrUser(receiver.ctx)
	if err != nil {
		return err
	}
	if len(strings.TrimSpace(email)) <= 0 {
		return exception.NewException(response.NotInputtedEmail)
	}

	currUser.UserEmail = email
	err = receiver.repo.SaveUser(currUser)
	return exception.ErrorHandle(err, response.DbExecuteError)
}
