package log

import (
	"bytes"
	"fmt"
	"github.com/sirupsen/logrus"
	"time"
)

// MyFormatter 自定义日志格式
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
