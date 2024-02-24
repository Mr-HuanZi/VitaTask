package hooks

import (
	"VitaTaskGo/internal/pkg/auth"
	"VitaTaskGo/internal/pkg/gateway"
	"github.com/sirupsen/logrus"
	"strconv"
)

func AuthUser(c *gateway.ChatClient, p gateway.Payload) {
	logrus.Debugf("AuthUser Hook: %v", p)
	if p.Data == nil {
		logrus.Errorln("AuthUser Hook: Data is nil")
		return
	}

	authorization := p.Data.(string)
	// 从请求头获取Token并解析
	claims, err := auth.ParseAuthorization(authorization)
	if err != nil {
		logrus.Errorln("Token解析失败：", err)
		return
	}
	logrus.Debugf("AuthUser Hook: %+v", claims)
	gateway.BingUserToClient(strconv.FormatUint(claims.UserId, 10), c.GetUniqueId())
}
