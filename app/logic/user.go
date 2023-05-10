package logic

import (
	"VitaTaskGo/app/constant"
	"VitaTaskGo/app/exception"
	"VitaTaskGo/app/extend"
	"VitaTaskGo/app/extend/user"
	"VitaTaskGo/app/model"
	"VitaTaskGo/app/response"
	"VitaTaskGo/app/types"
	"VitaTaskGo/library/db"
	"github.com/gin-gonic/gin"
	"github.com/gotidy/copy"
	"gorm.io/gorm"
	"strings"
	"time"
)

type UserLogic struct {
	Db  *gorm.DB
	ctx *gin.Context
}

func NewUserLogic(ctx *gin.Context) *UserLogic {
	return &UserLogic{
		Db:  db.Db, // 赋予ORM实例
		ctx: ctx,   // 传递上下文
	}
}

// CurrUser 返回当前登录用户
func (receiver UserLogic) CurrUser() (*model.User, error) {
	return user.CurrUser(receiver.ctx)
}

// UserExist 用户是否存在
func (receiver UserLogic) UserExist(uid uint64) bool {
	return user.Exist(uid)
}

// StoreSelf 保存用户信息
func (receiver UserLogic) StoreSelf(data types.UserInfoDto) (model.User, error) {
	// 获取登录用户uid
	uid, _ := receiver.ctx.Get(constant.CurrUidKey)
	if !receiver.UserExist(uid.(uint64)) {
		return model.User{}, exception.NewException(response.UserNotFound)
	}
	// 拷贝变量
	userData := model.User{}
	copiers := copy.New(func(c *copy.Options) {
		c.Skip = true
	})
	copiers.Copy(&userData, &data)

	err := receiver.Db.Model(&model.User{}).Where("id = ?", uid).Updates(userData).Error

	// 重新获取用户信息
	receiver.Db.Model(&model.User{}).First(&userData, uid)
	return userData, err
}

// ChangeAvatar 变更头像
func (receiver UserLogic) ChangeAvatar(avatar types.FileDto) error {
	currUser, err := user.CurrUser(receiver.ctx)
	if err != nil {
		return err
	}
	// 是否提供url
	if len(strings.TrimSpace(avatar.Url)) <= 0 {
		return exception.NewException(response.AvatarNotUploaded)
	}

	currUser.Avatar = avatar.Url
	err = receiver.Db.Save(&currUser).Error
	return exception.ErrorHandle(err, response.DbExecuteError)
}

// ChangePassword 修改当前用户密码
func (receiver UserLogic) ChangePassword(data types.ChangePasswordDto) error {
	currUser, err := user.CurrUser(receiver.ctx)
	if err != nil {
		return err
	}

	// 旧密码匹配
	if extend.Encryption(data.OldPassword) != currUser.UserPass {
		return exception.NewException(response.PassError)
	}

	currUser.UserPass = extend.Encryption(data.Password)
	// 记录密码修改的时间
	currUser.LastEditPass = time.Now().UnixMilli()
	err = receiver.Db.Save(&currUser).Error
	return exception.ErrorHandle(err, response.DbExecuteError)
}

// ChangeMobile 变更手机号
func (receiver UserLogic) ChangeMobile(mobile string) error {
	currUser, err := user.CurrUser(receiver.ctx)
	if err != nil {
		return err
	}
	if len(strings.TrimSpace(mobile)) <= 0 {
		return exception.NewException(response.NotInputtedMobile)
	}

	currUser.Mobile = mobile
	err = receiver.Db.Save(&currUser).Error
	return exception.ErrorHandle(err, response.DbExecuteError)
}

// ChangeEmail 变更电子邮箱地址
func (receiver UserLogic) ChangeEmail(email string) error {
	currUser, err := user.CurrUser(receiver.ctx)
	if err != nil {
		return err
	}
	if len(strings.TrimSpace(email)) <= 0 {
		return exception.NewException(response.NotInputtedEmail)
	}

	currUser.UserEmail = email
	err = receiver.Db.Save(&currUser).Error
	return exception.ErrorHandle(err, response.DbExecuteError)
}
