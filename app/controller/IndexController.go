package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

type IndexController struct {
}

func (cc *IndexController) Hello(c *gin.Context) {
	logrus.Infoln("Hello控制器执行了")
	c.String(http.StatusOK, "Hello World!")
}
