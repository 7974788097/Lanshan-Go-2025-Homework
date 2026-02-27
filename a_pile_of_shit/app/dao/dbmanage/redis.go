package dao

import (
	"a_pile_of_shit/app/model"
	"a_pile_of_shit/app/other"
	"fmt"
)

func getInformationOfRedis(Info *model.DBInfo) {
	port := other.GetenvDefault("REDIS_PORT", "6379")
	host := other.GetenvDefault("REDIS_HOST", "localhost")
	password := other.GetenvDefault("REDIS_PASSWORD", "123456")

	addr := fmt.Sprint(host, ":", port)
	Info.RedisAddr = addr
	Info.RedisPassword = password
}
