package service

import (
	"VitaTaskGo/internal/api/data"
	"VitaTaskGo/internal/api/model/dto"
	"VitaTaskGo/internal/pkg"
	"VitaTaskGo/internal/pkg/auth"
	"VitaTaskGo/internal/pkg/constant"
	"VitaTaskGo/internal/pkg/ws"
	"VitaTaskGo/internal/repo"
	"VitaTaskGo/pkg/exception"
	"VitaTaskGo/pkg/response"
	"github.com/duke-git/lancet/v2/slice"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type DialogService struct {
	Db   *gorm.DB
	ctx  *gin.Context
	repo repo.DialogRepo
}

func NewDialogService(tx *gorm.DB, ctx *gin.Context) *DialogService {
	return &DialogService{
		Db:   tx,  // 赋予ORM实例
		ctx:  ctx, // 传递上下文
		repo: data.NewDialogRepo(tx, ctx),
	}
}

// SendText 发送文本消息
func (receiver DialogService) SendText(dto dto.DialogSendTextDto) (*repo.DialogMsg, error) {
	// 检查对话是否存在
	dialog, err := receiver.repo.GetDialog(dto.DialogId)
	if err != nil {
		return nil, err
	}

	// 获取当前用户
	currUser, err := auth.CurrUser(receiver.ctx)
	if err != nil {
		return nil, err
	}

	// 当前用户是否在对话内
	if !receiver.repo.InDialog(dialog.ID, currUser.ID) {
		return nil, exception.NewException(response.NotInDialog)
	}

	msg, err := receiver.SaveMsg(dialog.ID, 0, "text", dto.Content)
	if err != nil {
		return nil, err
	}

	return msg, receiver.SendToDialog(dialog.ID, msg)
}

func (receiver DialogService) SendToDialog(dialogId uint, msgData *repo.DialogMsg) error {
	// 检查对话是否存在
	_, err := receiver.repo.GetDialog(dialogId)
	if err != nil {
		return err
	}

	// 获取对话成员
	members, err := data.NewDialogUserRepo(receiver.Db, receiver.ctx).GetDialogUsers(dialogId)
	if err != nil {
		return err
	}

	// 获取当前用户
	currUser, err := auth.CurrUser(receiver.ctx)
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

// SaveMsg 保存消息到数据库
// 不校验dialogId与userId，请自行在外部校验
func (receiver DialogService) SaveMsg(dialogId uint, userId uint64, t string, content string) (*repo.DialogMsg, error) {
	if userId <= 0 {
		// 获取当前用户
		userInfo, err := auth.CurrUser(receiver.ctx)
		if err != nil {
			return nil, err
		}
		userId = userInfo.ID
	}
	msg := &repo.DialogMsg{
		DialogId: dialogId,
		UserId:   userId,
		Type:     t,
		Content:  content,
	}

	err := data.NewDialogMsgRepo(receiver.Db, receiver.ctx).CreateDialogMsg(msg)
	if err != nil {
		return nil, exception.ErrorHandle(err, response.DbExecuteError)
	}

	// 关联其它内容
	msg.Dialog, _ = receiver.repo.GetDialog(dialogId)
	msg.UserInfo, _ = auth.CurrUser(receiver.ctx)

	return msg, nil
}

// Create 创建对话
func (receiver DialogService) Create(name, t string, members []uint64) (*repo.Dialog, error) {
	// 确认对话类型是否正确
	if !receiver.CheckType(t) {
		return nil, exception.NewException(response.DialogTypeError)
	}
	// 获取当前登录人
	currUser, err := auth.CurrUser(receiver.ctx)
	if err != nil {
		return nil, err
	}
	// 用户数组去重
	members = pkg.SliceUnique(members)
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
	dialog := repo.Dialog{Name: name, Type: t}
	// 开始创建
	err = receiver.Db.Transaction(func(tx *gorm.DB) error {
		// 此处要创建新的Repo
		dialogRepo := data.NewDialogRepo(tx, receiver.ctx)
		dialogUserRepo := data.NewDialogUserRepo(tx, receiver.ctx)
		userRepo := data.NewUserRepo(tx, receiver.ctx)
		// 创建对话
		err := dialogRepo.CreateDialog(&dialog)
		if err != nil {
			return err
		}

		for _, uid := range members {
			// 判断用户是否存在
			if !userRepo.Exist(uid) {
				return exception.NewException(response.UserNotFound)
			}
			// 创建成员记录
			err := dialogUserRepo.CreateDialogUser(&repo.DialogUser{DialogId: dialog.ID, UserId: uid})
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
func (receiver DialogService) Join(dialogId uint, members []uint64) error {
	// 检查对话是否存在
	dialog, err := receiver.repo.GetDialog(dialogId)
	if err != nil {
		return err
	}

	// 获取当前登录人
	currUser, err := auth.CurrUser(receiver.ctx)
	if err != nil {
		return err
	}

	// 获取对话成员列表
	oldMembers, err := data.NewDialogUserRepo(receiver.Db, receiver.ctx).GetDialogUsers(dialogId)
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
		dialogUserRepo := data.NewDialogUserRepo(tx, receiver.ctx)
		for _, uid := range members {
			if receiver.repo.InDialog(dialog.ID, uid) {
				// 如果成员已在对话中
				return exception.NewException(response.IsInDialog)
			}
			err := dialogUserRepo.CreateDialogUser(&repo.DialogUser{DialogId: dialog.ID, UserId: uid})
			if err != nil {
				return err
			}
		}

		return nil
	})
	return exception.ErrorHandle(err, response.JoinDialogFail)
}

// Exit 退出对话
func (receiver DialogService) Exit(dialogId uint, members []uint64) error {
	// 实例化Repo
	dialogUserRepo := data.NewDialogUserRepo(receiver.Db, receiver.ctx)
	// 获取对话成员列表
	oldMembers, err := dialogUserRepo.GetDialogUsers(dialogId)
	if err != nil {
		return exception.NewException(response.DbQueryError)
	}
	// 至少需要保留1个成员
	if len(members)-len(oldMembers) >= 1 {
		return exception.NewException(response.DialogKeep1Member)
	}
	return data.NewDialogUserRepo(receiver.Db, receiver.ctx).DeleteDialogUser(dialogId, members)
}

// Delete 删除对话
func (receiver DialogService) Delete(dialogId uint) error {
	err := receiver.Db.Transaction(func(tx *gorm.DB) error {
		var txErr error
		// 此处要创建新的Repo
		dialogRepo := data.NewDialogRepo(tx, receiver.ctx)
		dialogUserRepo := data.NewDialogUserRepo(tx, receiver.ctx)
		dialogMsgRepo := data.NewDialogMsgRepo(tx, receiver.ctx)

		// 删除消息记录
		txErr = dialogMsgRepo.DeleteDialogMsg(dialogId)
		if txErr != nil {
			return txErr
		}

		// 删除对话成员
		txErr = dialogUserRepo.DeleteDialogAllUser(dialogId)
		if txErr != nil {
			return txErr
		}

		// 删除对话
		txErr = dialogRepo.DeleteDialog(dialogId)
		return txErr
	})
	return exception.ErrorHandle(err, response.DialogDeleteFail)
}

func (receiver DialogService) CheckType(t string) bool {
	return slice.Contain(constant.GetDialogTypes(), t)
}

// MsgList 消息列表(不分页)
func (receiver DialogService) MsgList(dialogId uint) ([]repo.DialogMsg, error) {
	// 对话是否存在
	dialog, err := receiver.repo.GetDialog(dialogId)
	if err != nil {
		return nil, err
	}

	// 获取当前用户
	currUser, err := auth.CurrUser(receiver.ctx)
	if err != nil {
		return nil, err
	}

	// 当前用户是否在对话内
	if !receiver.repo.InDialog(dialog.ID, currUser.ID) {
		return nil, exception.NewException(response.NotInDialog)
	}

	// 拉取对话消息
	msgList, err := data.NewDialogMsgRepo(receiver.Db, receiver.ctx).ListDialogMsg(dialogId)
	return msgList, err
}
