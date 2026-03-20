package userLoginManage

import (
	"user-service/model"

	"gorm.io/gorm"
)

type UsersLoginInfo struct {
	*model.UsersLoginInfo
}
type dbOperation struct {
	mysql *gorm.DB
}

func Make(info *map[string]interface{}) *UsersLoginInfo {
	newInfo := &UsersLoginInfo{}
	newInfo.Username = (*info)["username"].(string)
	if (*info)["password"] != nil {
		newInfo.Password = (*info)["password"].(string)
	}
	return newInfo
}
func newdbOperation(db *model.DBcontext) *dbOperation {
	return &dbOperation{
		db.Mysql,
	}
}
func (user *UsersLoginInfo) AddInfo(dbContext *model.DBcontext) error {
	var err error
	db := newdbOperation(dbContext)

	err = db.SetUserLoginInfoOfMysql(user)
	if err != nil {
		return err
	}

	return nil
}
func (user *UsersLoginInfo) GetInfo(dbContext *model.DBcontext) (bool, error) {
	var err error
	var bo bool
	db := newdbOperation(dbContext)

	bo, err = db.GetUserLoginInfoOfMysql(user)
	if err != nil {
		return false, err
	}
	return bo, nil
}
