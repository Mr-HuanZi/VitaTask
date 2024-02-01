package main

import (
	"VitaTaskGo/internal/api"
	"VitaTaskGo/internal/api/middleware"
	"VitaTaskGo/internal/pkg/workflow"
	"VitaTaskGo/pkg/config"
	"VitaTaskGo/pkg/db"
	"VitaTaskGo/pkg/log"
	"flag"
	"github.com/gin-gonic/gin"
	"strconv"
)

var configFile = flag.String("f", "config/app.yaml", "the config file")

func main() {
	flag.Parse()

	// 读取配置
	err := config.Load(*configFile)
	if err != nil {
		panic(err)
	}

	// 初始化日志
	log.InitLogsDriver()
	// 初始化数据库
	initDatabases()
	// 初始化工作流
	workflow.Init()
	// 初始化Gin
	r := gin.Default()
	// 注册中间件
	middleware.Register(r)
	// 注册路由
	api.Routers(r)
	// 注册WebSocket路由
	api.WebSocketRouters(r)
	// 绑定Host与Port
	_ = r.Run(config.Instances.App.Host + ":" + strconv.Itoa(config.Instances.App.Port))
}

// 初始化数据库
func initDatabases() {
	err := db.Init(db.DsnConfig{
		Drive:  "mysql",
		Host:   config.Instances.Mysql.Host,
		Port:   config.Instances.Mysql.Port,
		User:   config.Instances.Mysql.User,
		Pass:   config.Instances.Mysql.Password,
		Dbname: config.Instances.Mysql.DbName,
		Prefix: config.Instances.Mysql.Prefix,
	})

	if err != nil {
		panic("Database connection failed")
	}
}
