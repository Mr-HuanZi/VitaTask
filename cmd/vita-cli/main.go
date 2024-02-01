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
	log.InitLogsDriver()
	// 初始化数据库
	initDatabases()
	// 命令行处理
	cli.FlagHandle()
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
