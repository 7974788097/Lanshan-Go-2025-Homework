package userinfomanage

import (
	"Homework_LanShan/new_shit/app/model"
	"log"

	"gorm.io/gorm"
)

func AutoUserInfoManageDB(dbContext *gorm.DB) {
	err := dbContext.AutoMigrate(model.UsersInfo{})
	if err != nil {
		log.Fatalln(err)
	}
}
func (U *UserInfoStruct) addInMysql(dbContext *gorm.DB) error {
	result := dbContext.Create(U)
	return result.Error
}
