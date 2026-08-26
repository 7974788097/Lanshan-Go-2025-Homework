package messagemanage

import (
	"a_pile_of_shit/app/model"
	"errors"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func connectMessageMysql(info *model.DBInfo) {
	path := info.MysqlDataPath
	var err error
	for i := range 5 {
		dbContext.Mysql, err = gorm.Open(mysql.Open(path), &gorm.Config{})
		if err == nil {
			break
		}
		log.Printf("Mysql of message connect error %d times: %s", i, err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatal("Mysql of message connect error: ", err)
	}
	log.Println("Mysql of message connect success")
	err = dbContext.Mysql.AutoMigrate(&model.MessageInfo{})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Mysql of message creat success")
}

func SetMassageOfMysql(info *model.MessageInfo) error {
	var err error
	err = dbContext.Mysql.Model(model.MessageInfo{}).Create(info).Error

	return err
} //向mysql中添加文章，返回错误

func UpdateMessageOfMysql(info *model.MessageInfo) (bool, error) {
	result := dbContext.Mysql.Model(model.MessageInfo{}).Where("message_id = ?", info.MessageID).Updates(info)
	if result.Error != nil {
		return false, result.Error
	}
	if result.RowsAffected == 0 {
		return false, nil
	}
	return true, nil
} //更新mysql中文章的信息，，返回错误 xs

func GetMessageOfMysqlByID(info *MessageInfo) (bool, error) {
	result := dbContext.Mysql.Model(model.MessageInfo{}).Where("message_id = ?", info.MessageID).Find(&(info.Info))
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return true, result.Error
	}
	return true, nil
} //从mysql中获取文章信息以及是否存在，返回bool和error
func GetMessageOfMysqlByMessageName(info *MessageInfo) (bool, error) {
	result := dbContext.Mysql.Model(model.MessageInfo{}).Where("message_name LIKE ? AND status = ?", "%"+info.MessageName+"%", 1).Find(&(info.Info))
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return true, result.Error
	}
	return true, nil
}
func GetMessageOfMysqlByUsername(info *MessageInfo) (bool, error) {
	result := dbContext.Mysql.Model(model.MessageInfo{}).Where("username = ?", info.Username).Find(&(info.Info))
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return true, result.Error
	}
	return true, nil
}
func DeleteMessageOfMysql(info *model.MessageInfo) error {
	ID := info.MessageID
	result := dbContext.Mysql.Model(model.MessageInfo{}).Where("message_id = ?", ID).Delete(info)
	return result.Error
}
