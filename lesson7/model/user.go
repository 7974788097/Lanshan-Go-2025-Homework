package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username string `gorm:"size:18,not null,comment:账号"`
	Password string `gorm:"size:18,default:123456,not null,comment:密码"`
}
