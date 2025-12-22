package dao

import (
	"context"
	"log"
	"moddle/model"
	"os"

	"github.com/redis/go-redis/v9"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func getInformationAboutDatabase() *model.DatabaseInformation {
	file, err := os.Open("./config/config.yaml")
	if err != nil {
		panic(err)
	}
	defer func() { _ = file.Close() }()
	var cfg model.DatabaseInformation
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&cfg)
	if err != nil {
		panic(err)
	}
	return &cfg
}
func connectMySQL(MySQL string) {
	var err error
	db, err = gorm.Open(mysql.Open(MySQL), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	log.Println("自动迁移成功")
	err = db.AutoMigrate(&model.User{})
	if err != nil {
		panic(err)
	}
}

var cli *redis.Client

func connectRedis(addr string, password string) {
	ctx := context.Background()
	cli = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
	})
	_, err := cli.Ping(ctx).Result()
	if err != nil {
		panic(err)
	}
}
