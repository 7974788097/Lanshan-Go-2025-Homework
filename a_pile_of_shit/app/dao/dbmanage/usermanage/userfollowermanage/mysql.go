package userfollowermanage

import (
	"a_pile_of_shit/app/model"
	"errors"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func connectUserFollowerMysql(info *model.DBInfo) {
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
	err = dbContext.Mysql.AutoMigrate(&model.UsersFollowerInfo{})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Mysql of user Mute creat success")
}
func SetUserFollowerInfoOfMysql(user *model.UsersFollowerInfo) error {
	var err error
	result := dbContext.Mysql.Create(user)
	err = result.Error
	return err
} //向mysql中添加user，返回错误
func UpdateUserFollowerInfoOfMysql(user *model.UsersFollowerInfo) (bool, error) {
	err := errors.New("this module is discontinued")
	return false, err
} //不使用，仅为使接口成立
func DeleteUserFollowerInfoOfMysql(user *model.UsersFollowerInfo) error {
	var err error
	for range 5 {
		err = dbContext.Mysql.Where("follow_username = ? AND followed_username = ?", user.FollowUsername, user.FollowedUsername).Limit(1).Delete(&model.UsersFollowerInfo{}).Error
		if err == nil || errors.Is(err, gorm.ErrRecordNotFound) {
			break
		}
	}
	return err
} //从mysql中软删除对应user

func GetUserFollowerInfoOfMysql(user *model.UsersFollowerInfo) ([]string, bool, error) {
	followerName := user.FollowUsername

	var total int64
	dbContext.Mysql.Model(&model.UsersFollowerInfo{}).Where("follow_username = ?", followerName).Count(&total)
	if total == 0 {
		return nil, false, nil
	}
	pageSize := 20
	page := (int(total) + pageSize - 1) / pageSize
	var temporary []model.UsersFollowerInfo
	var end []string
	for i := range page {
		offset := pageSize * i
		result := dbContext.Mysql.Where("follow_username = ?", followerName).Limit(pageSize).Offset(offset).Find(&temporary)
		if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, false, result.Error
		}
		for _, j := range temporary {
			end = append(end, j.FollowedUsername)
		}
	}
	return end, true, nil
} //从mysql中获取用户关注的username并更改user信息以及是否存在，返回userInfo,bool和error
func GetUserFollowedInfoOfMysql(user *model.UsersFollowerInfo) ([]string, bool, error) {
	followedName := user.FollowedUsername

	var total int64
	dbContext.Mysql.Model(&model.UsersFollowerInfo{}).Where("followed_username = ?", followedName).Count(&total)
	if total == 0 {
		return nil, false, nil
	}
	pageSize := 20
	page := (int(total) + pageSize - 1) / pageSize
	var temporary []model.UsersFollowerInfo
	var end []string
	for i := range page {
		offset := pageSize * i
		result := dbContext.Mysql.Where("followed_username = ?", followedName).Limit(pageSize).Offset(offset).Find(&temporary)
		if result.Error != nil && !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, false, result.Error
		}
		for _, j := range temporary {
			end = append(end, j.FollowUsername)
		}
	}
	return end, true, nil
} // //从mysql中获取用户的粉丝的username并更改user信息以及是否存在，返回userInfo,bool和error
