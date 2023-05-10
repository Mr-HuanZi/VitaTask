package logic

import (
	"VitaTaskGo/app/exception"
	"VitaTaskGo/app/extend"
	"VitaTaskGo/app/extend/user"
	"VitaTaskGo/app/model"
	"VitaTaskGo/app/response"
	"VitaTaskGo/app/types"
	"VitaTaskGo/library/config"
	"VitaTaskGo/library/db"
	"github.com/gin-gonic/gin"
	"time"
)

type MemberLogic struct {
	ctx *gin.Context
}

func NewMemberLogic(ctx *gin.Context) *MemberLogic {
	return &MemberLogic{
		ctx: ctx, // 传递上下文
	}
}

// Lists 成员列表
func (receiver MemberLogic) Lists(query types.MemberListsQuery) (types.PagedResult[model.User], error) {
	var (
		members []model.User
		count   int64
		where   = make(map[string]interface{})
	)
	/* 查询 Start */
	if query.Id > 0 {
		// 如果是ID查询就只允许一个条件
		where["id"] = query.Id
	} else {
		if query.Status > 0 {
			where["user_status"] = query.Id
		}
	}
	tx := db.Db.Model(model.User{}).Where(where)
	// Map查询无法进行LIKE，所以下面直接使用Where字符串查询
	if query.Username != "" {
		tx = tx.Where("user_login LIKE ?", "%"+query.Username+"%")
	}
	if query.Nickname != "" {
		tx = tx.Where("user_nickname LIKE ?", "%"+query.Nickname+"%")
	}
	if query.Mobile != "" {
		tx = tx.Where("mobile LIKE ?", "%"+query.Mobile+"%")
	}
	if query.Email != "" {
		tx = tx.Where("user_email LIKE ?", "%"+query.Email+"%")
	}
	/* 查询 End */
	if err := tx.Count(&count).Error; err != nil {
		return types.PagedResult[model.User]{
			Items: nil,
			Total: 0,
			Page:  1,
		}, err
	}

	if err := tx.Scopes(db.Paginate(&query.Page, &query.PageSize)).Find(&members).Error; err != nil {
		return types.PagedResult[model.User]{
			Items: nil,
			Total: 0,
			Page:  1,
		}, err
	}
	return types.PagedResult[model.User]{
		Items: members,
		Total: count,
		Page:  int64(query.Page),
	}, nil
}

// SimpleList 简单成员列表
func (receiver MemberLogic) SimpleList(key string) []types.SimpleMemberList {
	var simpleList []types.SimpleMemberList
	tx := db.Db.Model(model.User{}).Where("user_status = ?", 1)
	if key != "" {
		tx = tx.Where("user_login LIKE ? OR user_nickname LIKE ?", "%"+key+"%", "%"+key+"%")
	}
	tx.Find(&simpleList)
	return simpleList
}

// Create 创建成员
func (receiver MemberLogic) Create(data types.MemberCreate) (*model.User, error) {
	var (
		count int64
		u     = new(model.User)
	)

	// 校验密码格式
	if !extend.PassFormat(data.Password) {
		return nil, exception.NewException(response.RegPassFormatError)
	}
	// 验证用户名是否存在
	if err := db.Db.Model(model.User{}).Where("user_login = ?", data.Username).Count(&count).Error; err != nil {
		return nil, err
	}
	u.UserNickname = data.Nickname
	u.UserLogin = data.Username
	u.UserPass = extend.Encryption(data.Password)
	u.UserEmail = data.Email
	u.Mobile = data.Mobile
	u.UserStatus = 1 // 启用
	u.LastEditPass = time.Now().Unix()

	// 创建用户
	if err := db.Db.Create(&u).Error; err != nil {
		return nil, err
	}
	return u, nil
}

// ChangeUserStatus 更改用户状态
func (receiver MemberLogic) ChangeUserStatus(uid uint64, status int) error {
	// 验证用户名是否存在
	one, err := receiver.GetOne(uid)
	if err != nil {
		return err
	}
	// 设置为禁用状态
	return db.Db.Model(&one).Updates(map[string]interface{}{"user_status": status}).Error
}

// ResetPassword 重置用户密码
func (receiver MemberLogic) ResetPassword(uid uint64) error {
	one, err := receiver.GetOne(uid)
	if err != nil {
		return err
	}
	// 校验密码格式
	if !extend.PassFormat(config.Instances.Member.DefaultPass) {
		return exception.NewException(response.RegPassFormatError)
	}
	// 生成新密码
	userPass := extend.Encryption(config.Instances.Member.DefaultPass)
	// 修改密码的时间
	lastEditPass := time.Now().Unix()
	// 保存新密码
	return db.Db.Model(&one).Updates(model.User{UserPass: userPass, LastEditPass: lastEditPass}).Error
}

// GetOne 获取一个用户
func (receiver MemberLogic) GetOne(uid uint64) (*model.User, error) {
	var u *model.User
	err := db.Db.First(&u, uid).Error
	if err != nil {
		return nil, err
	}

	// 如果UID小于0
	if u.ID <= 0 {
		return nil, exception.NewException(response.UserNotFound)
	}
	return u, nil
}

// ChangeSuper 改变一个成员的超级管理员状态
// 请在外部确认 super 值是否合法
func (receiver MemberLogic) ChangeSuper(uid uint64, super int8) error {
	// 获取当前用户
	currUser, err := user.CurrUser(receiver.ctx)
	if err != nil {
		return err
	}
	if !user.IsSuper(currUser) {
		return exception.NewException(response.CurrUserNotSuper)
	}
	// 用户是否存在
	if !user.Exist(uid) {
		return exception.NewException(response.UserNotFound)
	}
	// 不能是自己
	if currUser.ID == uid {
		return exception.NewException(response.UserSuperChangeSelf)
	}

	err = db.Db.Model(&model.User{}).Where("id = ?", uid).Update("super", super).Error
	if err != nil {
		return exception.ErrorHandle(err, response.DbExecuteError)
	}
	return nil
}
