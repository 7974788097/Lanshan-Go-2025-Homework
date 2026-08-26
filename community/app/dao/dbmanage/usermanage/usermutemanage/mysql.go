package usermutemanage

import (
	"a_pile_of_shit/app/model"
	"errors"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func connectUserMuteMysql(info *model.DBInfo) {
	path := info.MysqlUserPath
	var err error
	for i := range 5 {
		dbContext.Mysql, err = gorm.Open(mysql.Open(path), &gorm.Config{})
		if err == nil {
			break
		}
		log.Printf("mysql of user Mute connect error %d times: %s", i, err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatal("Mysql of user Mute connect error: ", err)
	}
	log.Println("Mysql of user Mute connect success")
	err = dbContext.Mysql.AutoMigrate(&UsersMuteInfo{})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Mysql of user Mute creat success")
}
func SetUserMuteInfoOfMysql(user *UsersMuteInfo) error {
	var err error
	result := dbContext.Mysql.Create(user)
	err = result.Error

	return err
} //向mysql中添加user，返回错误
func UpdateUserMuteInfoOfMysql(user *UsersMuteInfo) (bool, error) {
	result := dbContext.Mysql.Where("username = ?", user.Username).Updates(user)
	if result.Error != nil {
		return false, result.Error
	}
	if result.RowsAffected == 0 {
		return false, nil
	}
	return true, nil
} //更新mysql中user的值，返回错误
func DeleteUserMuteInfoOfMysql(user *UsersMuteInfo) error {
	var err error
	for range 5 {
		err = dbContext.Mysql.Where("username = ?", user.Username).Delete(user).Error
		if err == nil {
			break
		}
	}
	return err
} //从mysql中软删除对应user

func GetUserMuteInfoOfMysql(user *UsersMuteInfo) (UsersMuteInfo, bool, error) {
	username := user.Username
	result := dbContext.Mysql.Where("username = ?", username).First(user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return UsersMuteInfo{}, false, nil
		}
		return UsersMuteInfo{}, true, result.Error
	}
	return *user, true, nil
} //从mysql中获取并更改user信息以及是否存在，返回userInfo,bool和error
