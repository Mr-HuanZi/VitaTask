package command

import (
	"VitaTaskGo/app/model"
	"VitaTaskGo/library/db"
	"flag"
	"github.com/sirupsen/logrus"
	"os"
	"strings"
)

var (
	test = ""
)

func init() {
	Register("migrate", AutoMigrate)
}

func AutoMigrate(f *flag.FlagSet) bool {
	f.StringVar(&test, "test", "", "测试")
	// 忽略错误
	_ = f.Parse(os.Args[2:])
	if len(strings.TrimSpace(test)) > 0 {
		logrus.Debugln("命令行测试", test)
	}
	// 执行数据迁移
	err := db.Db.Set("gorm:table_options", "ENGINE=InnoDB").
		AutoMigrate(&model.Dialog{}, &model.DialogMsg{}, &model.DialogUser{})
	if err != nil {
		logrus.Errorln(err)
		return false
	}
	return false
}
