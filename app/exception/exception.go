package exception

import (
	"fmt"
)

type Exception struct {
	Code    int
	Message string
}

func (receiver Exception) Error() string {
	return receiver.Message
}

func (receiver Exception) GetCode() int {
	return receiver.Code
}

// NewException 实例化异常
func NewException(code int, args ...interface{}) error {
	var message string
	if len(args) > 2 {
		message = fmt.Sprintf(args[0].(string), args...)
	} else if len(args) > 0 {
		message = fmt.Sprintf("%s", args...)
	} else {
		// 此处是 args参数为空的情况，消息交给response包处理
		message = ""
	}
	return &Exception{
		Code:    code,
		Message: message,
	}
}
