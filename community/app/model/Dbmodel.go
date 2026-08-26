package model

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type DBcontext struct {
	Mysql *gorm.DB
	Redis *redis.Client
}
