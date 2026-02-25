package userloginmanage

import (
	"a_pile_of_shit/app/model"
	"a_pile_of_shit/app/other"
	"context"
	"errors"
	"log"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

func connectUserLoginRedis(info *model.DBInfo) {
	ctx := context.Background()
	addr := info.RedisAddr
	password := info.RedisPassword
	dbContext.Redis = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
	})
	var err error
	for i := range 5 {
		_, err = dbContext.Redis.Ping(ctx).Result()
		if err == nil {
			break
		}
		log.Printf("redis of use Login connect error %d times: %s", i, err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatal("redis of use Login connect error: ", err)
	}
	log.Println("redis of use Login connect success")
}
func SetUserOfRedis(user *UsersLoginInfo, cacheTime time.Duration) error {
	ctx := context.Background()
	info := make(map[string]interface{})

	info["password"] = user.Password
	info["salt"] = user.Salt
	info["permission"] = user.Permission
	err := other.RedisHSetEX(dbContext.Redis, ctx, cacheTime, user.Username, info)
	return err
} //向redis中设置用户密码，盐值
func GetUserOfRedis(user *UsersLoginInfo) (bool, error) {

	ctx := context.Background()
	info, err := dbContext.Redis.HGetAll(ctx, user.Username).Result()

	if err != nil {
		if errors.Is(err, redis.Nil) {
			user.Password = ""
			user.Salt = ""
			return false, nil
		}
		return true, err
	}
	user.Password = info["password"]
	user.Salt = info["salt"]
	permission, _ := strconv.Atoi(info["permission"])
	user.Permission = uint8(permission)
	if user.Salt == "" || user.Password == "" {
		return false, nil
	}
	return true, nil
} //从redis中获取并更改user的密码，盐值，并返回用户是否存在

// ...................以下是登录认证相关.............................

func LogInRedis(username string) error {
	ctx := context.Background()
	info := make(map[string]interface{})
	info["loginState"] = true
	err := other.RedisHSetEX(dbContext.Redis, ctx, model.CacheTime, username, info)
	return err
}
func LogOutRedis(username string) error {
	ctx := context.Background()
	info := make(map[string]interface{})
	info["loginState"] = false
	err := other.RedisHSetEX(dbContext.Redis, ctx, model.CacheTime, username, info)
	return err
}
func GetLoginStateRedis(username string) (bool, error) {
	ctx := context.Background()
	var bo bool
	info, err := dbContext.Redis.HGet(ctx, username, "loginState").Result()
	if errors.Is(redis.Nil, err) {
		err = nil
		info = "0"
	}
	if info == "1" {
		bo = true
	} else {
		bo = false
	}
	return bo, err
}
