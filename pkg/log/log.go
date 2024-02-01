package log

import (
	"VitaTaskGo/pkg/path_x"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"path"
	"time"
)

var logBaseDir = "./logs" // 日志基础路径

// InitLogsDriver 初始化日志
func InitLogsDriver(logFilePath string, ginLogFilePath string) error {
	// 创建日志目录
	err := os.MkdirAll(logBaseDir, 0755)
	if err != nil {
		panic(err)
	}
	// 日志等级
	logrus.SetLevel(logrus.DebugLevel)
	// 输出格式
	logrus.SetFormatter(&MyFormatter{TimestampFormat: time.DateTime})

	// 拼接完整的日志路径
	logFilePath = path.Join(logBaseDir, logFilePath)
	ginLogFilePath = path.Join(logBaseDir, ginLogFilePath)
	// 检查传入的路径格式是否合法
	if !path_x.PathValid(logFilePath) || !path_x.PathValid(ginLogFilePath) {
		return errors.New("log file path is invalid")
	}

	// 分割日志
	logger := &lumberjack.Logger{
		Filename:   logFilePath,
		MaxSize:    1,  // 一个文件最大为1M
		MaxAge:     30, // 一个文件最多同时存在30天
		MaxBackups: 5,  // 最多同时保存5份文件
		LocalTime:  true,
		Compress:   false,
	}
	fileWriter := io.MultiWriter(logger)
	logrus.SetOutput(fileWriter)

	logrus.Infoln("logrus initialized")

	// Gin日志重定向
	InitGinLogDriver(ginLogFilePath)

	return nil
}

// InitGinLogDriver Gin日志重定向
func InitGinLogDriver(filePath string) {
	logger := &lumberjack.Logger{
		Filename:   filePath,
		MaxSize:    1,  // 一个文件最大为1M
		MaxAge:     30, // 一个文件最多同时存在30天
		MaxBackups: 5,  // 最多同时保存5份文件
		LocalTime:  true,
		Compress:   false,
	}
	gin.DisableConsoleColor()                             // 不需要颜色
	gin.DefaultWriter = io.MultiWriter(os.Stdout, logger) // 同时输出到屏幕
}
