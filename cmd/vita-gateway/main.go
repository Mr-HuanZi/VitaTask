package main

import (
	"VitaTaskGo/internal/gateway"
	_ "VitaTaskGo/internal/gateway/hooks"
	"VitaTaskGo/pkg/config"
	"VitaTaskGo/pkg/log"
	"flag"
	"github.com/gin-gonic/gin"
	"strconv"
)

var configFile = flag.String("f", "config/app.yaml", "the config file")

func main() {
	initialize()

	// 初始化Gin
	r := gin.Default()
	// 注册路由
	gateway.Routers(r)
	// 注册WebSocket路由
	gateway.WebSocketRouters(r)
	// 绑定Host与Port
	_ = r.Run(config.Get().Gateway.Host + ":" + strconv.Itoa(config.Get().Gateway.Port))
}

func initialize() {
	flag.Parse()
	// 读取配置
	err := config.Load(*configFile)
	if err != nil {
		panic(err)
	}
	// 初始化日志
	logErr := log.InitLogsDriver("gateway.log", "gateway-gin.log")
	if logErr != nil {
		panic(logErr)
	}
}
