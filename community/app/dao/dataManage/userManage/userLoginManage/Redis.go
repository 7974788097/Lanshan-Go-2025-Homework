package userloginmmanage

import (
	"Homework_LanShan/new_shit/app/model"
	"Homework_LanShan/new_shit/app/other"
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

func (U *UserLoginInfoStruct) addInRedis(dbContext *redis.Client, time time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), model.TimeOutLimit)
	defer cancel()
	info := make(map[string]any)
	info["permission"] = U.Permission
	err := other.RedisHSetEX(dbContext, ctx, time, "user:info"+U.Username, info)
	return err
}
