package main

import (
	"VitaTaskGo/internal/cli"
	"VitaTaskGo/pkg/config"
	"VitaTaskGo/pkg/db"
	"VitaTaskGo/pkg/log"
)

func main() {
	// 读取配置
	err := config.Load("./app.yaml")
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
	// 命令行处理
	cli.FlagHandle()
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
