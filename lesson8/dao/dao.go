package dao

import (
	"context"
	"errors"
	"log"
	"moddle/model"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var db *gorm.DB

func ConnectDatabase() {
	cfg := getInformationAboutDatabase()
	connectMySQL(cfg.Mysql.Path)
	connectRedis(cfg.Redis.Addr, cfg.Redis.Password)
}
func DeleteUserFromRedis(username string) {
	ctx := context.Background()
	cli.Del(ctx, username)
}
func AddUserIntoMySQL(username string, password string) {
	add := &model.User{Username: username, Password: password}
	db.Create(add)
}
func AddUserIntoRedis(username string, password string, time time.Duration) {
	ctx := context.Background()
	DeleteUserFromRedis(username)
	cli.SetEx(ctx, username, password, time)
}
func AddUser(user *model.User) {
	AddUserIntoMySQL(user.Username, user.Password)
	AddUserIntoRedis(user.Username, user.Password, 20*time.Second)
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
func FindUserFromRedis(username string) error {
	ctx := context.Background()
	name, err := cli.Get(ctx, username).Result()
	if err != nil {
		AddUserIntoRedis(username, "", 5*time.Minute)
		return model.NoUser
	}
	if name == "" {
		return model.RepetionLogin
	}
	return nil
}

func FindUserFromDatabase(username string) bool {
	err := FindUserFromRedis(username)
	if err == nil {
		return true
	}
	if errors.Is(err, model.RepetionLogin) {
		return false
	}
	if FindUserFromMySQL(username) {
		password := GetPasswordFromMySQl(username)
		AddUserIntoRedis(username, password, 20*time.Second)
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
	if err == nil {
		return password, true
	}
	if !errors.Is(err, redis.Nil) {
		log.Println(err)
	}
	return "", false
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
	updatePasswordFromMySQL(username, password)
	updatePasswordFromRedis(username, password)
}
func DelectUserFromMySQL(username string) bool {
	var user model.User
	result := db.Where("username = ?", username).Delete(&user)
	if result.RowsAffected == 0 {
		return false
	}
	return true
}
func DelectUserFromRedis(username string) {
	ctx := context.Background()
	cli.SetEx(ctx, username, "", 20*time.Minute)
}
func DelectUserFromDatabase(username string) bool {
	bo := DelectUserFromMySQL(username)
	DelectUserFromRedis(username)
	return bo
}
