package model

import "gorm.io/gorm"

type CommentInfo struct {
	gorm.Model   `json:"-"`
	CommentID    int64   `gorm:"index;not null;unique"`
	Content      string  `gorm:"type:text;not null;comment:文章正文"`
	Username     string  `gorm:"index;size:54;not null;comment:账号"`
	ChildNodeID  []int64 `gorm:"serializer:json;type:json;not null;comment:回复评论"`
	ParentNodeID int64   `gorm:"not null;comment:所回复的评论或文章"`
	Path         string  `gorm:"index;type:varchar(512);not null;comment:路径"`
}
type MessageInfo struct {
	gorm.Model  `json:"-"`
	MessageID   int64  `gorm:"index;not null;unique"`
	MessageName string `gorm:"index;size:90;not null;comment:文章标题" json:"MessageName,omitempty"`
	Content     string `gorm:"type:mediumtext;not null;comment:文章正文" json:"Content,omitempty"`
	Username    string `gorm:"index;size:54;not null;comment:发布用户名" json:"Username,omitempty"`
	Permission  uint8  `gorm:"not null;default:1;comment:阅读所需的权限等级" json:"Permission,omitempty"`
	Status      uint8  `gorm:"not null;default:0;comment:1-草稿 2-已发布" json:"Status,omitempty"`
}
