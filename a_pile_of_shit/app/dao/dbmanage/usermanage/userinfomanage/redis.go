package userinfomanage

import (
	"a_pile_of_shit/app/model"
	"a_pile_of_shit/app/other"
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

func connectUserInfoRedis(info *model.DBInfo) {
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
		log.Printf("redis of userinfo connect error %d times: %s", i, err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatal("redis of userinfo connect error: ", err)
	}
	log.Println("redis of userinfo connect success")
}
func setUserInfoOfRedis(user *UsersInfo, cacheTime time.Duration) error {
	ctx := context.Background()
	info := make(map[string]interface{})

	info["permission"] = user.Permission

	err := other.RedisHSetEX(dbContext.Redis, ctx, cacheTime, user.Username, info)
	return err
} //向redis中设置用户密码,权限，盐值等所有用户信息
