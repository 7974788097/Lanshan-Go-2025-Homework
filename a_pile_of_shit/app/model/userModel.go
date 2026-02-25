package model

import "gorm.io/gorm"

// 更改后得手动更新usermanage中redis的相关操作
type UsersLoginInfo struct {
	gorm.Model `json:"-"`
	Username   string `gorm:"index;primaryKey;size:54;not null;unique;comment:账号"`
	Password   string `gorm:"size:64;not null;comment:密码"`
	Salt       string `gorm:"size:64" json:"-"`
	Permission uint8  `gorm:"default:1;not null;comment:权限"`
}
type UsersInfo struct {
	gorm.Model `json:"-"`
	Username   string `gorm:"index;primaryKey;size:54;not null;unique;comment:账号"`
	Permission uint8  `gorm:"default:1;not null;comment:权限"`
}
type UsersMuteInfo struct {
	gorm.Model `json:"-"`
	Username   string `gorm:"index;primaryKey;size:54;not null;unique;comment:账号"`
	MuteLabel  [2]int `gorm:"type:json;not null;serializer:json;comment:禁言标签 [0]==1为未禁言,[0]==2为禁言,[1]为结束的时间戳"`
	MuteReason string `gorm:"type:text;not null;comment:禁言原因"`
}
type UsersFollowerInfo struct {
	gorm.Model       `json:"-"`
	FollowUsername   string `gorm:"index;size:54;not null;comment:发起关注的用户名"`
	FollowedUsername string `gorm:"index;size:54;not null;comment:被关注的用户名"`
}
