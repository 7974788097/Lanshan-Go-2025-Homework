package commentmanage

import (
	"a_pile_of_shit/app/model"
	"errors"
	"log"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func connectCommentMysql(info *model.DBInfo) {
	path := info.MysqlDataPath
	var err error
	for i := range 5 {
		dbContext.Mysql, err = gorm.Open(mysql.Open(path), &gorm.Config{})
		if err == nil {
			break
		}
		log.Printf("Mysql of comment connect error %d times: %s", i, err)
		time.Sleep(2 * time.Second)
	}
	if err != nil {
		log.Fatal("Mysql of comment connect error: ", err)
	}
	log.Println("Mysql of comment connect success")
	err = dbContext.Mysql.AutoMigrate(&model.CommentInfo{})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Mysql of comment creat success")
}

func SetCommentOfMysql(info *model.CommentInfo) error {
	var err error
	for range 5 {
		err = dbContext.Mysql.Create(info).Error
		if err == nil {
			break
		}
		time.Sleep(1 * time.Second)
	}
	return err
} //向mysql中添加文章，返回错误

func UpdateCommentOfMysql(info *model.CommentInfo) (bool, error) {
	result := dbContext.Mysql.Model(model.CommentInfo{}).Where("comment_id = ?", info.CommentID).Updates(info)
	if result.Error != nil {
		return false, result.Error
	}
	if result.RowsAffected == 0 {
		return false, nil
	}
	return true, nil
} //更新mysql中评论的信息，，返回错误

func GetCommentOfMysqlByPath(info *CommentInfo) (bool, error) {
	var getInfo []model.CommentInfo
	result := dbContext.Mysql.Model(model.CommentInfo{}).Where("path LIKE ?", info.Path+"%").Find(&getInfo)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return true, result.Error
	}
	info.Info = getInfo
	return true, nil
} //从mysql中获取评论信息以及是否存在，获取的信息储存在Info中，返回bool和error
func GetCommentOfMysqlByID(info *CommentInfo) (bool, error) {
	var getInfo []model.CommentInfo
	result := dbContext.Mysql.Model(model.CommentInfo{}).Where("comment_id = ?", info.CommentID).Find(&getInfo)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return true, result.Error
	}
	info.Info = getInfo
	return true, nil
} //从mysql中获取评论信息以及是否存在，获取的信息储存在Info中，返回bool和error
func GetCommentOfMysqlByUsername(info *CommentInfo) (bool, error) {
	var getInfo []model.CommentInfo
	result := dbContext.Mysql.Model(model.CommentInfo{}).Where("username = ?", info.Username).Find(&getInfo)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return false, nil
		}
		return true, result.Error
	}
	info.Info = getInfo
	return true, nil
} //从mysql中获取评论信息以及是否存在，获取的信息储存在Info中，返回bool和error
func DeleteCommentOfMysqlByPath(info *model.CommentInfo) error {
	result := dbContext.Mysql.Model(model.CommentInfo{}).Where("path LIKE ?", info.Path+"%").Delete(model.CommentInfo{})
	return result.Error
}
