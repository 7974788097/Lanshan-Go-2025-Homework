package model

import "gorm.io/gorm"

type UsersLoginInfo struct {
	gorm.Model `json:"-"`
	Username   string `gorm:"index;primaryKey;size:54;not null;unique;comment:账号"`
	Password   string `gorm:"size:64;not null;comment:密码"`
	Salt       string `gorm:"size:64" json:"-"`
	Permission uint8  `gorm:"default:1;not null;comment:权限"`
}
