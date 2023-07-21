package service

import (
	"VitaTaskGo/internal/biz"
	"VitaTaskGo/internal/data"
	"VitaTaskGo/internal/dto"
	"VitaTaskGo/internal/pkg"
	"VitaTaskGo/internal/pkg/auth"
	"VitaTaskGo/internal/pkg/config"
	"VitaTaskGo/internal/pkg/exception"
	"VitaTaskGo/internal/pkg/response"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"time"
)

type MemberService struct {
	Db   *gorm.DB
	ctx  *gin.Context
	repo biz.UserRepo
}

func NewMemberService(tx *gorm.DB, ctx *gin.Context) *MemberService {
	return &MemberService{
		Db:   tx,  // 赋予ORM实例
		ctx:  ctx, // 传递上下文
		repo: data.NewUserRepo(tx, ctx),
	}
}

// Lists 成员列表
func (receiver MemberService) Lists(query dto.MemberListsQuery) (dto.PagedResult[biz.User], error) {
	list, total, err := receiver.repo.PageListUser(query)
	return dto.PagedResult[biz.User]{
		Items: list,
		Total: total,
		Page:  int64(query.Page),
	}, exception.ErrorHandle(err, response.DbQueryError)
}

// SimpleList 简单成员列表
func (receiver MemberService) SimpleList(key string) []dto.SimpleMemberList {
	return receiver.repo.SimpleList(key)
}

// Create 创建成员
func (receiver MemberService) Create(data dto.MemberCreate) (*biz.User, error) {
	u := new(biz.User)

	// 校验密码格式
	if !pkg.PassFormat(data.Password) {
		return nil, exception.NewException(response.RegPassFormatError)
	}
	// 验证用户名是否存在
	if receiver.repo.ExistByUsername(data.Username) {
		return nil, exception.NewException(response.RegUsernameExists)
	}

	// 给各个字段赋值
	u.UserNickname = data.Nickname
	u.UserLogin = data.Username
	u.UserPass = pkg.Encryption(data.Password)
	u.UserEmail = data.Email
	u.Mobile = data.Mobile
	u.UserStatus = 1 // 启用
	u.LastEditPass = time.Now().Unix()

	// 创建用户
	err := receiver.repo.CreateUser(u)
	return u, err
}

// ChangeUserStatus 更改用户状态
func (receiver MemberService) ChangeUserStatus(uid uint64, status int) error {
	// 验证用户名是否存在
	if !receiver.repo.Exist(uid) {
		return exception.NewException(response.UserNotFound)
	}
	// 设置为禁用状态
	return receiver.repo.UpdateUserStatus(uid, status)
}

// ResetPassword 重置用户密码
func (receiver MemberService) ResetPassword(uid uint64) error {
	// 验证用户名是否存在
	if !receiver.repo.Exist(uid) {
		return exception.NewException(response.UserNotFound)
	}
	// 校验密码格式
	if !pkg.PassFormat(config.Instances.Member.DefaultPass) {
		return exception.NewException(response.RegPassFormatError)
	}
	// 生成新密码
	userPass := pkg.Encryption(config.Instances.Member.DefaultPass)
	// 保存新密码
	return receiver.repo.UpdateUserPass(uid, userPass)
}

// ChangeSuper 改变一个成员的超级管理员状态
// 请在外部确认 super 值是否合法
func (receiver MemberService) ChangeSuper(uid uint64, super int8) error {
	// 获取当前用户
	currUser, err := auth.CurrUser(receiver.ctx)
	if err != nil {
		return err
	}

	if !auth.IsSuper(currUser) {
		return exception.NewException(response.CurrUserNotSuper)
	}

	// 用户是否存在
	if !receiver.repo.Exist(uid) {
		return exception.NewException(response.UserNotFound)
	}

	// 不能是自己
	if currUser.ID == uid {
		return exception.NewException(response.UserSuperChangeSelf)
	}

	// 执行更新
	return exception.ErrorHandle(receiver.repo.UpdateUserSuper(uid, super), response.DbExecuteError)
}
