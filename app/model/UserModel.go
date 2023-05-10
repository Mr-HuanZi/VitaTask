package model

import (
	"VitaTaskGo/app/extend/time_tool"
	"VitaTaskGo/library/config"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID                uint64 `json:"id,omitempty" gorm:"primaryKey"`
	Sex               int8   `json:"sex"`
	Birthday          string `json:"birthday" gorm:"default:null"`
	LastLoginTime     uint64 `json:"lastLoginTime" gorm:"default:null"`
	LastLoginIp       string `json:"lastLoginIp"`
	CreateTime        uint64 `json:"createTime" gorm:"autoUpdateTime"`
	UpdateTime        uint64 `json:"updateTime" gorm:"autoUpdateTime"`
	UserStatus        uint8  `json:"userStatus"`
	UserLogin         string `json:"userLogin" gorm:"index"`
	UserPass          string `json:"-"`
	UserNickname      string `json:"userNickname"`
	UserEmail         string `json:"userEmail"`
	Avatar            string `json:"avatar"`
	Signature         string `json:"signature"`
	UserActivationKey string `json:"userActivationKey"`
	Mobile            string `json:"mobile"`
	LockTime          int64  `json:"lockTime,omitempty"`
	ErrorSum          uint8  `json:"errorSum,omitempty"`
	First             uint8  `json:"first,omitempty"`
	LastEditPass      int64  `json:"lastEditPass"`
	Openid            string `json:"openid"`
	Super             int8   `json:"super"`
}

func (receiver *User) TableName() string {
	return config.Instances.Mysql.Prefix + "user"
}

func (receiver *User) AfterFind(*gorm.DB) (err error) {
	// 转换日期格式
	if receiver.Birthday != "" {
		receiver.Birthday, _ = time_tool.ChangeFormat(time.RFC3339, time.DateOnly, receiver.Birthday)
	}
	return
}
