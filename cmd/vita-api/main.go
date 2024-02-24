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
	err = log.InitLogsDriver("app.log", "gin.log")
	if err != nil {
		panic(err)
	}
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
	// 绑定Host与Port
	_ = r.Run(config.Get().App.Host + ":" + strconv.Itoa(config.Get().App.Port))
}

// 初始化数据库
func initDatabases() {
	err := db.Init(db.DsnConfig{
		Drive:  "mysql",
		Host:   config.Get().Mysql.Host,
		Port:   config.Get().Mysql.Port,
		User:   config.Get().Mysql.User,
		Pass:   config.Get().Mysql.Password,
		Dbname: config.Get().Mysql.DbName,
		Prefix: config.Get().Mysql.Prefix,
	})

	if err != nil {
		panic("Database connection failed")
	}
}
