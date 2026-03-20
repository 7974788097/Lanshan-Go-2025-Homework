package userLoginManage

import (
	"errors"
	"fmt"
	"log"
	"user-service/model"

	"gorm.io/gorm"
)

func AutoUserLoginInfoDB(dbContext *gorm.DB) {
	err := dbContext.AutoMigrate(&model.UsersLoginInfo{})
	if err != nil {
		log.Fatal(fmt.Errorf("AutoUserLoginInfoDB: %w", err))
	}
	log.Println("Mysql of use Login creat success")
}
func (dbContext *dbOperation) SetUserLoginInfoOfMysql(user *UsersLoginInfo) error {
	result := dbContext.mysql.Create(user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrDuplicatedKey) || result.RowsAffected == 0 {
			return result.Error
		}
		return result.Error
	}

	return nil
}
func (dbContext *dbOperation) GetUserLoginInfoOfMysql(user *UsersLoginInfo) (bool, error) {
	username := user.Username
	result := dbContext.mysql.Where("username = ?", username).First(user)
	if result.Error != nil {
		fmt.Println("mysql err:", result.Error)
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return true, result.Error
	}
	return true, nil
}
