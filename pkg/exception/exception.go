package exception

import (
	"fmt"
	"strings"
)

type Exception struct {
	Code    int
	Message string
}

// MultipleException 允许将多个错误合并为一个Error的异常Type
type MultipleException struct {
	Code    int
	Message string
	Errors  []error
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

func NewMultipleException(code int, args ...interface{}) *MultipleException {
	var message string
	if len(args) > 2 {
		message = fmt.Sprintf(args[0].(string), args...)
	} else if len(args) > 0 {
		message = fmt.Sprintf("%s", args...)
	} else {
		// 此处是 args参数为空的情况，消息交给response包处理
		message = ""
	}
	return &MultipleException{
		Code:    code,
		Message: message,
		Errors:  nil,
	}
}

func (r *MultipleException) Error() string {
	resMsg := r.Message
	if r.Errors != nil && len(r.Errors) > 0 {
		errorsMsg := make([]string, len(r.Errors))
		for i, err := range r.Errors {
			errorsMsg[i] = err.Error()
		}
		resMsg = resMsg + "\nMore:\n" + strings.Join(errorsMsg, "\n")
	}
	return resMsg
}

func (r *MultipleException) GetCode() int {
	return r.Code
}

// AddErr 追加错误信息
func (r *MultipleException) AddErr(err error) *MultipleException {
	if err == nil {
		return r
	}

	if r.Errors == nil {
		r.Errors = []error{err}
	} else {
		r.Errors = append(r.Errors, err)
	}

	return r
}
