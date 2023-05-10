package logic

import (
	"VitaTaskGo/app/constant"
	"VitaTaskGo/app/exception"
	"VitaTaskGo/app/extend"
	"VitaTaskGo/app/extend/user"
	"VitaTaskGo/app/model"
	"VitaTaskGo/app/modules/ws"
	"VitaTaskGo/app/response"
	"VitaTaskGo/app/types"
	"VitaTaskGo/library/db"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type DialogLogic struct {
	Db  *gorm.DB
	ctx *gin.Context
}

func NewDialogLogic(ctx *gin.Context) *DialogLogic {
	return &DialogLogic{
		Db:  db.Db, // 赋予ORM实例
		ctx: ctx,   // 传递上下文
	}
}

// GetOne 获取单条记录
func (receiver DialogLogic) GetOne(dialogId uint) (*model.Dialog, error) {
	var dialog *model.Dialog
	err := receiver.Db.Model(&model.Dialog{}).First(&dialog, dialogId).Error
	return dialog, exception.ErrorHandle(err, response.DialogNotExist)
}

// SendText 发送文本消息
func (receiver DialogLogic) SendText(dto types.DialogSendTextDto) (*model.DialogMsg, error) {
	// 检查对话是否存在
	dialog, err := receiver.GetOne(dto.DialogId)
	if err != nil {
		return nil, err
	}
	// 当前用户是否在对话内
	if !receiver.In(dialog.ID, 0) {
		return nil, exception.NewException(response.NotInDialog)
	}

	msg, err := receiver.SaveMsg(dialog.ID, 0, "text", dto.Content)
	if err != nil {
		return nil, err
	}

	return msg, receiver.SendToDialog(dialog.ID, msg)
}

func (receiver DialogLogic) SendToDialog(dialogId uint, msgData *model.DialogMsg) error {
	// 检查对话是否存在
	_, err := receiver.GetOne(dialogId)
	if err != nil {
		return err
	}
	// 获取对话成员
	members, err := receiver.GetMembers(dialogId)
	if err != nil {
		return err
	}
	// 获取当前用户
	currUser, err := user.CurrUser(receiver.ctx)
	if err != nil {
		return err
	}
	// 将消息推送给对话成员
	for _, member := range members {
		// 跳过自己
		if member.UserId == currUser.ID {
			continue
		}
		// 只推送给在线的成员
		if ws.Online(member.UserId) {
			err := ws.Send(ws.GetClient(member.UserId), "chat", msgData)
			if err != nil {
				return err
			}
			logrus.Debugf("已将消息发送给UID[%d]", member.UserId)
		}
	}

	return nil
}

// In 用户是否在对话内
// userId 传递0值时取当前登录用户
func (receiver DialogLogic) In(dialogId uint, userId uint64) bool {
	var count int64
	if userId <= 0 {
		// 获取当前用户,忽略错误
		userInfo, _ := user.CurrUser(receiver.ctx)
		userId = userInfo.ID
	}
	// 用户是否在对话内
	receiver.Db.Model(&model.DialogUser{}).
		Where("dialog_id = ?", dialogId).
		Where("user_id = ?", userId).
		Count(&count)

	return count > 0
}

// SaveMsg 保存消息到数据库
// 不校验dialogId与userId，请自行在外部校验
func (receiver DialogLogic) SaveMsg(dialogId uint, userId uint64, t string, content string) (*model.DialogMsg, error) {
	if userId <= 0 {
		// 获取当前用户
		userInfo, err := user.CurrUser(receiver.ctx)
		if err != nil {
			return nil, err
		}
		userId = userInfo.ID
	}
	msg := &model.DialogMsg{
		DialogId: dialogId,
		UserId:   userId,
		Type:     t,
		Content:  content,
	}

	err := receiver.Db.Create(msg).Error
	if err != nil {
		return nil, exception.ErrorHandle(err, response.DbExecuteError)
	}

	// 关联其它内容
	msg.Dialog, _ = receiver.GetOne(dialogId)
	msg.UserInfo, _ = user.CurrUser(receiver.ctx)

	return msg, nil
}

// Create 创建对话
func (receiver DialogLogic) Create(name, t string, members []uint64) (*model.Dialog, error) {
	// 确认对话类型是否正确
	if !receiver.CheckType(t) {
		return nil, exception.NewException(response.DialogTypeError)
	}
	// 获取当前登录人
	currUser, err := user.CurrUser(receiver.ctx)
	if err != nil {
		return nil, err
	}
	// 用户数组去重
	members = extend.SliceUnique(members)
	// 加入成员
	if len(members) <= 0 {
		return nil, exception.NewException(response.DialogMemberEmpty)
	}
	// 如果是 C2C 类型
	if t == constant.DialogTypeC2C {
		// 只允许一个聊天对象
		if len(members) > 1 {
			return nil, exception.NewException(response.DialogC2COvercrowding)
		}
		// 另一个聊天对象是不是自己
		if currUser.ID == members[0] {
			return nil, exception.NewException(response.DialogMemberIsMe)
		}
	}

	// 创建对话实例
	dialog := model.Dialog{Name: name, Type: t}
	// 开始创建
	err = receiver.Db.Transaction(func(tx *gorm.DB) error {
		// 创建对话
		err := tx.Create(&dialog).Error
		if err != nil {
			return err
		}

		for _, uid := range members {
			// 判断用户是否存在
			userLogic := NewUserLogic(receiver.ctx)
			userLogic.Db = tx
			if !userLogic.UserExist(uid) {
				return exception.NewException(response.UserNotFound)
			}
			// 创建成员记录
			err := tx.Create(&model.DialogUser{DialogId: dialog.ID, UserId: uid}).Error
			if err != nil {
				return err
			}
		}
		return nil
	})
	return &dialog, exception.ErrorHandle(err, response.DialogCreateFail)
}

// Join 加入对话
// 如果成员已在对话中，会返回一个错误
func (receiver DialogLogic) Join(dialogId uint, members []uint64) error {
	// 检查对话是否存在
	dialog, err := receiver.GetOne(dialogId)
	if err != nil {
		return err
	}
	// 获取当前登录人
	currUser, err := user.CurrUser(receiver.ctx)
	if err != nil {
		return err
	}
	// 获取对话成员列表
	oldMembers, err := receiver.GetMembers(dialogId)
	if err != nil {
		return exception.NewException(response.DbQueryError)
	}
	if dialog.Type == constant.DialogTypeC2C {
		// c2c对话超员
		if len(oldMembers) >= 2 {
			return exception.NewException(response.DialogC2COvercrowding)
		}
		// 只允许一个聊天对象
		if len(members) > 1 {
			return exception.NewException(response.DialogC2COvercrowding)
		}
		// 另一个聊天对象是不是自己
		if currUser.ID == members[0] {
			return exception.NewException(response.DialogMemberIsMe)
		}
	}
	// 开始插入
	err = receiver.Db.Transaction(func(tx *gorm.DB) error {
		for _, uid := range members {
			if receiver.In(dialog.ID, uid) {
				// 如果成员已在对话中
				return exception.NewException(response.IsInDialog)
			}
			err := tx.Create(&model.DialogUser{DialogId: dialog.ID, UserId: uid}).Error
			if err != nil {
				return err
			}
		}

		return nil
	})
	return exception.ErrorHandle(err, response.JoinDialogFail)
}

// Exit 退出对话
func (receiver DialogLogic) Exit(dialogId uint, members []uint64) error {
	// 获取对话成员列表
	oldMembers, err := receiver.GetMembers(dialogId)
	if err != nil {
		return exception.NewException(response.DbQueryError)
	}
	// 至少需要保留1个成员
	if len(members)-len(oldMembers) >= 1 {
		return exception.NewException(response.DialogKeep1Member)
	}
	return receiver.Db.Where("dialog_id = ?", dialogId).Where("user_id IN ?", members).Delete(&model.DialogUser{}).Error
}

// Delete 删除对话
func (receiver DialogLogic) Delete(dialogId uint) error {
	err := receiver.Db.Transaction(func(tx *gorm.DB) error {
		var txErr error
		// 删除消息记录
		txErr = tx.Where("dialog_id = ?", dialogId).Delete(&model.DialogMsg{}).Error
		if txErr != nil {
			return txErr
		}

		// 删除对话成员
		txErr = tx.Where("dialog_id = ?", dialogId).Delete(&model.DialogUser{}).Error
		if txErr != nil {
			return txErr
		}

		// 删除对话
		txErr = tx.Delete(&model.Dialog{}, dialogId).Error
		return txErr
	})
	return exception.ErrorHandle(err, response.DialogDeleteFail)
}

func (receiver DialogLogic) CheckType(t string) bool {
	return slice.Contain(constant.GetDialogTypes(), t)
}

// GetMembers 获取对话成员
func (receiver DialogLogic) GetMembers(dialogId uint) ([]model.DialogUser, error) {
	var (
		members []model.DialogUser
	)

	err := receiver.Db.Model(&model.DialogUser{}).Where("dialog_id = ?", dialogId).Find(&members).Error
	return members, err
}

// MsgList 消息列表(不分页)
func (receiver DialogLogic) MsgList(dialogId uint) ([]model.DialogMsg, error) {
	var msgList []model.DialogMsg
	// 对话是否存在
	dialog, err := receiver.GetOne(dialogId)
	if err != nil {
		return nil, err
	}
	// 当前用户是否在对话内
	if !receiver.In(dialog.ID, 0) {
		return nil, exception.NewException(response.NotInDialog)
	}
	// 拉取对话消息
	err = receiver.Db.Model(&model.DialogMsg{}).
		Preload("Dialog").
		Preload("UserInfo").
		Where("dialog_id = ?", dialogId).
		Find(&msgList).Error
	return msgList, err
}

// GetMessagePreload 带预加载的获取单条 消息 记录方法
func (receiver DialogLogic) GetMessagePreload(msgId uint64) (*model.DialogMsg, error) {
	var msg *model.DialogMsg
	err := receiver.Db.Model(&model.DialogMsg{}).
		Preload("Dialog").
		Preload("UserInfo").
		First(&msg, msgId).Error
	return msg, exception.ErrorHandle(err, response.DialogNotExist)
}
