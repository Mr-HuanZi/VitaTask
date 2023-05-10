package log

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"os"
	"time"
)

type MyFormatter struct {
	TimestampFormat string
}

func (receiver *MyFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	// 默认时间格式
	if receiver.TimestampFormat == "" {
		receiver.TimestampFormat = time.DateTime
	}

	entry.HasCaller()

	logString := fmt.Sprintf("[%s] [%s] \n%v", entry.Time.Format(receiver.TimestampFormat), entry.Level, entry.Message)
	b.WriteString(logString)
	b.WriteByte('\n')
	return b.Bytes(), nil
}

// InitLogsDriver 初始化日志
func InitLogsDriver() {
	// 创建日志目录
	err := os.MkdirAll("./logs", os.ModePerm)
	if err != nil {
		panic(err)
	}
	// 日志等级
	logrus.SetLevel(logrus.DebugLevel)
	// 输出格式
	logrus.SetFormatter(&MyFormatter{TimestampFormat: time.DateTime})
	// 日志文件名
	logFileName := "./logs/app.log"
	// 输出路径
	//logFile, _ := os.OpenFile("./logs/app.log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	// 分割日志
	logger := &lumberjack.Logger{
		Filename:   logFileName,
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
	InitGinLogDriver()
}

// InitGinLogDriver Gin日志重定向
func InitGinLogDriver() {
	logFile, _ := os.OpenFile("./logs/gin.log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	gin.DisableConsoleColor()                              // 不需要颜色
	gin.DefaultWriter = io.MultiWriter(os.Stdout, logFile) // 同时输出到屏幕
}
