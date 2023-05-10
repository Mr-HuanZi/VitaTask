package exception

import "github.com/sirupsen/logrus"

// ErrorHandle 错误处理
// 根据err变量决定返回值，如果err == nil那么返回nil
// 如果err != nil那么返回 Exception 并且记录错误日志
// args 参数为错误日志前缀说明
func ErrorHandle(err error, code int, args ...string) error {
	if err != nil {
		if e, ok := err.(*Exception); ok {
			// 如果是自定义错误则直接返回
			return e
		}
		// 如果错误不为空，记录日志并且返回异常信息
		logrus.Errorln(args, err)
		return NewException(code)
	}
	return nil
}
