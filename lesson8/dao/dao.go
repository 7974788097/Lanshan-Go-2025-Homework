package dao

import (
	"context"
	"errors"
	"moddle/model"
	"time"

	"gorm.io/gorm"
)

var db *gorm.DB

func ConnectDatabase() {
	cfg := getInformation()
	connectMySQL(cfg.Mysql.Path)
	connectRedis(cfg.Redis.Addr, cfg.Redis.Password)
}

func AddUserIntoMySQL(username string, password string) {
	add := &model.User{Username: username, Password: password}
	db.Create(add)
}
func AddUserIntoRedis(username string, password string) {
	ctx := context.Background()
	cli.SetEx(ctx, username, password, 20*time.Minute)
}
func AddUser(user *model.User) {
	AddUserIntoMySQL(user.Username, user.Password)
	AddUserIntoRedis(user.Username, user.Password)
}

// 判断用户名是否存在
func FindUserFromMySQL(username string) bool {
	var user model.User
	result := db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		if user.Username == "" || errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false
		}
		panic(result.Error)
	}
	return true
}
func FindUserFromRedis(username string) bool {
	ctx := context.Background()
	_, err := cli.Get(ctx, username).Result()
	if err != nil {
		return false
	}
	return true
}

func FindUserFromDatabase(username string) bool {
	if FindUserFromRedis(username) {
		return true
	}
	if FindUserFromMySQL(username) {
		password := GetPasswordFromMySQl(username)
		AddUserIntoRedis(username, password)
		return true
	}
	return false
}

// 判断密码是否正确
func JudgePassword(username string, input string) bool {
	password := GetPasswordFromDatabase(username)
	if password == input {
		return true
	}
	return false
}

// 获取原本的密码
func GetPasswordFromMySQl(username string) string {
	var user model.User
	db.Where("username = ?", username).First(&user)
	return user.Password
}
func GetPasswordFromRedis(username string) (string, bool) {
	ctx := context.Background()
	password, err := cli.Get(ctx, username).Result()
	var bo bool
	if err != nil {
		bo = true
	} else {
		bo = false
	}
	return password, bo
}
func GetPasswordFromDatabase(username string) string {
	password, exist := GetPasswordFromRedis(username)
	if exist {
		return password
	}
	return GetPasswordFromMySQl(username)
}

// 更新密码
func updatePasswordFromMySQL(username string, password string) {
	var user model.User
	db.Where("username = ?", username).First(&user)
	db.Model(&user).Update("password", password)
}
func updatePasswordFromRedis(username string, password string) {
	ctx := context.Background()
	cli.SetEx(ctx, username, password, 20*time.Minute)
}
func UpdatePassword(username string, password string) {
	updatePasswordFromRedis(username, password)
	updatePasswordFromMySQL(username, password)
}
