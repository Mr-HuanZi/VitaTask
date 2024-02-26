package services

import (
	"VitaTaskGo/internal/pkg/gateway"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ChatService struct {
	ctx *gin.Context
}

func NewChatService(ctx *gin.Context) *ChatService {
	return &ChatService{ctx: ctx}
}

func (receiver *ChatService) SendToUser(user string, msg string) error {
	return receiver.SendToUsers([]string{user}, msg)
}

func (receiver *ChatService) SendToUsers(users []string, msg string) error {
	// 遍历users
	for _, userid := range users {
		// 获取该用户的Client ID
		clientId := gateway.GetClientID(userid)
		if clientId == "" {
			logrus.Warnf("gateway error: 用户[%s]不存在或未登录", userid)
			continue
		}
		// 获取该用户的Client
		client := gateway.GetClient(clientId)
		// 发送消息
		if client != nil {
			client.Send([]byte(msg))
		} else {
			logrus.Warnf("gateway error: 找不到用户[%s]Client[%s]", userid, clientId)
		}
	}
	return nil
}
