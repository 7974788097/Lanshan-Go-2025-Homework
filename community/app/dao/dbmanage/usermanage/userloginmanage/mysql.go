package userloginmanage

import (
	"a_pile_of_shit/app/model"
	"errors"
	"fmt"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func connectUserLoginMysql(info *model.DBInfo) {
	path := info.MysqlUserPath
	var err error
	for i := range 5 {
		dbContext.Mysql, err = gorm.Open(mysql.Open(path), &gorm.Config{})
		if err == nil {
			break
		}
		log.Printf("mysql of use Login connect error %d times: %s", i, err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatal("Mysql of use Login connect error: ", err)
	}
	log.Println("Mysql of use Login connect success")
	err = dbContext.Mysql.AutoMigrate(&UsersLoginInfo{})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Mysql of use Login creat success")
}
func SetUserLoginOfMysql(user *UsersLoginInfo) error {
	var err error
	result := dbContext.Mysql.Create(user)
	err = result.Error

	return err
} //向mysql中添加user，返回错误
func UpdateUserLoginOfMysql(user *UsersLoginInfo) (bool, error) {
	result := dbContext.Mysql.Where("username = ?", user.Username).Updates(user)
	if result.Error != nil {
		return false, result.Error
	}
	if result.RowsAffected == 0 {
		return false, nil
	}
	return true, nil
} //更新mysql中user的值，返回错误
func DeleteUserLoginOfMysql(user *UsersLoginInfo) error {
	var err error
	for range 5 {
		err = dbContext.Mysql.Delete(user).Error
		if err == nil {
			break
		}
	}
	return err
} //从mysql中软删除对应user

func GetUserLoginOfMysql(user *UsersLoginInfo) (UsersLoginInfo, bool, error) {
	username := user.Username
	result := dbContext.Mysql.Where("username = ?", username).First(user)
	if result.Error != nil {
		fmt.Println("mysql err:", result.Error)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return UsersLoginInfo{}, false, nil
		}
		return UsersLoginInfo{}, true, result.Error
	}
	return *user, true, nil
} //从mysql中获取并更改user信息以及是否存在，返回userInfo,bool和error
