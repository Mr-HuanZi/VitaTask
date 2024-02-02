package db

import (
	"VitaTaskGo/pkg/config"
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"strconv"
	"time"
)

var (
	Db *gorm.DB
)

type DsnConfig struct {
	Drive  string
	Host   string
	Port   int
	User   string
	Pass   string
	Dbname string
	Prefix string
}

// GormWriter sql监听
type GormWriter struct {
}

// Printf 实现gorm/logger.Writer接口
func (m *GormWriter) Printf(format string, v ...interface{}) {
	// 记录日志
	logrus.Infof(format, v...)
}

func Init(dsnConfig DsnConfig) error {
	var (
		openErr error
		db      *gorm.DB
	)

	// user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		dsnConfig.User,
		dsnConfig.Pass,
		dsnConfig.Host,
		strconv.Itoa(dsnConfig.Port),
		dsnConfig.Dbname,
	)

	// 调试模式下启用SQL日志
	if config.Get().App.Debug {
		newLogger := logger.New(
			&GormWriter{},
			logger.Config{
				LogLevel:                  logger.Info, // 日志级别
				IgnoreRecordNotFoundError: true,        // 忽略ErrRecordNotFound（记录未找到）错误
				Colorful:                  false,       // 禁用彩色打印
			},
		)
		db, openErr = gorm.Open(mysql.Open(dsn), &gorm.Config{
			Logger: newLogger,
		})
	} else {
		db, openErr = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	}

	if openErr != nil {
		return openErr
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	// 设置连接可复用的最大时间
	// 服务器设置了10分钟超时，那么这里就设置9分钟
	sqlDB.SetConnMaxLifetime(time.Minute * 9)

	Db = db // 赋值给包变量
	return nil
}

// Paginate 分页
func Paginate(page, pageSize *int) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if *page <= 0 {
			*page = 1
		}
		if *pageSize < 10 {
			*pageSize = 10
		}
		offset := (*page - 1) * *pageSize
		return db.Offset(offset).Limit(*pageSize)
	}
}
