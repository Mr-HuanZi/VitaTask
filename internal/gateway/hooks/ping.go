package hooks

import (
	"VitaTaskGo/internal/pkg/gateway"
	"github.com/sirupsen/logrus"
	"time"
)

func PingHandle(c *gateway.ChatClient, p gateway.Payload) {
	// 收到ping消息后返回消息给客户端
	err := c.SendMsg("ping", time.Now().Unix())
	if err != nil {
		logrus.Errorln("PingHandle SendMsg Error: ", err)
	}
}
