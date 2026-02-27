package model

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type DbError struct {
	MysqlError error
	RedisError error
}
type DBInfo struct {
	MysqlInfo
	RedisInfo
}
type MysqlInfo struct {
	MysqlUserPath string
	MysqlDataPath string
}
type RedisInfo struct {
	RedisAddr     string
	RedisPassword string
}
type DBcontext struct {
	Mysql *gorm.DB
	Redis *redis.Client
}
