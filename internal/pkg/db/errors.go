package db

import (
	"VitaTaskGo/internal/pkg/exception"
	"VitaTaskGo/internal/pkg/response"
	"errors"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// FirstQueryErrorHandle 数据库查询单条记录的错误处理
func FirstQueryErrorHandle(err error, code int) error {
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return exception.NewException(code)
		}

		logrus.Errorln("数据库查询出错", err.Error())
		return exception.NewException(response.DbQueryError)
	}
	return nil
}
