package model

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type DBInfo struct {
	MysqlInfo
	RedisInfo
}
type MysqlInfo struct {
	MysqlPath string
}
type RedisInfo struct {
	RedisAddr     string
	RedisPassword string
}
type DBcontext struct {
	Mysql *gorm.DB
	Redis *redis.Client
}
