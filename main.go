package main

import (
	"VitaTaskGo/app/command"
	"VitaTaskGo/app/middleware"
	"VitaTaskGo/library/config"
	"VitaTaskGo/library/db"
	"VitaTaskGo/library/log"
	"VitaTaskGo/routers"
	"github.com/gin-gonic/gin"
	"strconv"
)

func main() {
	// 读取配置
	err := config.Load("./app.yml")
	if err != nil {
		panic(err)
	}

	// 初始化日志
	log.InitLogsDriver()
	// 初始化数据库
	initDatabases()
	// 命令行处理
	if !command.FlagHandle() {
		return
	}
	// 初始化Gin
	r := gin.Default()
	// 注册中间件
	middleware.Register(r)
	// 注册路由
	routers.ApiRouters(r)
	// 注册WebSocket路由
	routers.WebSocketRouters(r)
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
