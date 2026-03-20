package dao

import (
	"log"
	"user-service/dao/userLoginManage"
	"user-service/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func AutoUserManage(mysql *gorm.DB) {
	userLoginManage.AutoUserLoginInfoDB(mysql)
}
func GetDBInfo() *model.DBInfo {
	var mysqlInfo = model.MysqlInfo{
		"root:123456@tcp(127.0.0.1:3306)/test?charset=utf8",
	}
	var redisInfo = model.RedisInfo{
		"localhost:6379",
		"123456",
	}
	return &model.DBInfo{
		mysqlInfo,
		redisInfo,
	}
}
func ConnectUserDB(dbInfo *model.DBInfo) *model.DBcontext {
	dbContext := model.DBcontext{}
	var err error
	dbContext.Mysql, err = gorm.Open(mysql.Open(dbInfo.MysqlPath), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	return &dbContext
}

type UserDB struct {
	*model.DBcontext
}

func NewUserDB(dbContext *model.DBcontext) *UserDB {
	return &UserDB{
		dbContext,
	}
}

type UserOperation interface {
	AddInfo(dbContext *model.DBcontext) error
	GetInfo(dbContext *model.DBcontext) (bool, error)
	//省略其他方法
}

func Add(Info UserOperation, dbContext *model.DBcontext) error {
	return Info.AddInfo(dbContext)
}
func Get(Info UserOperation, dbContext *model.DBcontext) (bool, error) {
	return Info.GetInfo(dbContext)
}
