package middleware

import (
	"VitaTaskGo/internal/pkg/auth"
	"VitaTaskGo/internal/pkg/constant"
	"VitaTaskGo/pkg/config"
	"VitaTaskGo/pkg/response"
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
	"strings"
)

// Register 注册全局中间件
func Register(r *gin.Engine) {
	if config.Instances.App.Debug {
		//r.Use(RequestRecord())
	}
}

// CheckLogin 验证用户是否登录
func CheckLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
		var authorization string
		if strings.TrimSpace(c.GetHeader("Upgrade")) == "websocket" {
			// websocket请求
			authorization = c.GetHeader("Sec-WebSocket-Protocol")
		} else {
			// 普通Http请求
			authorization = c.GetHeader("Authorization")
		}
		// 从请求头获取Token并解析
		claims, err := auth.ParseAuthorization(authorization)
		if err != nil {
			logrus.Errorln("Token解析失败：", err)
			c.JSON(http.StatusUnauthorized, response.Exception(response.SignatureMissing))
			c.Abort()
			return
		}
		// 将user信息保存到上下文
		c.Set(constant.CurrUidKey, claims.UserId)
	}
}

// RequestRecord 请求记录器
func RequestRecord() gin.HandlerFunc {
	return func(c *gin.Context) {
		body, err := io.ReadAll(c.Request.Body)
		if err != nil {
			logrus.Errorln("读取请求体失败，错误信息为：", err)
			return
		}
		// 需要把Body再次放回去，不然就会出现EOF的情况
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		header := ""
		for k, v := range c.Request.Header {
			header += fmt.Sprintf("%s: %s\n", k, v)
		}
		logrus.Infof("URI:%s\nHeader\n%s\nBody\n%s", c.Request.RequestURI, header, body)
	}
}
