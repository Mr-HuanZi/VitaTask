package response

import (
	"VitaTaskGo/app/exception"
	"VitaTaskGo/app/validator"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

var (
	successMessage = "操作成功"
	failMessage    = "操作失败"
)

type MessageBody struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Code    int         `json:"code"`
	custom  bool        // 是否自定义
}

// Custom 自定义返回
func Custom(message string, code int, data interface{}) MessageBody {
	return MessageBody{
		Data:    data,
		Message: message,
		Code:    code,
	}
}

func Success() MessageBody {
	return MessageBody{
		Data:    nil,
		Message: successMessage,
		Code:    http.StatusOK,
	}
}

func SuccessData(data interface{}) MessageBody {
	return MessageBody{
		Data:    data,
		Message: successMessage,
		Code:    http.StatusOK,
	}
}

// Auto 自动判断是返回错误还是成功
func Auto(data interface{}, err error) MessageBody {
	if err != nil {
		return Error(err)
	} else {
		return SuccessData(data)
	}
}

// Error 统一错误处理
// 非自定义错误时，返回默认消息
// 如果不是自定义默认消息的，用这个方法，否则请使用 FormatError
func Error(entity error) MessageBody {
	message := failMessage
	code := SystemFail
	custom := false
	// 如果是自定义异常
	if err, ok := entity.(*exception.Exception); ok {
		// 传递过来的消息为空，直接从状态码中获取
		if len(strings.TrimSpace(err.Message)) <= 0 {
			message = GetMessage(err.Code)
		} else {
			message = err.Message
		}
		code = err.Code
		custom = true
	} else {
		// 非自定义错误则打印日志
		logrus.Errorln(entity)
	}
	return MessageBody{
		Data:    nil,
		Message: message,
		Code:    code,
		custom:  custom,
	}
}

// HandleFormVerificationFailed 表单验证失败专用的错误返回
func HandleFormVerificationFailed(err error) MessageBody {
	return MessageBody{
		Data:    nil,
		Message: validator.FailHandle(err),
		Code:    FormVerificationFailed,
	}
}

// Exception 自定义异常通用处理
func Exception(code int) MessageBody {
	return Error(exception.NewException(code))
}
