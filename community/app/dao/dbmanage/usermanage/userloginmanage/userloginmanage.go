package userloginmanage

import (
	"a_pile_of_shit/app/model"
	"a_pile_of_shit/app/other"
	"log"
	"time"
)

var dbContext *model.DBcontext

func ConnectUserLoginDB(info *model.DBInfo) {
	dbContext = &model.DBcontext{}
	connectUserLoginMysql(info)
	connectUserLoginRedis(info)
}

type UsersLoginInfo struct {
	model.UsersLoginInfo
}

func Make(info *map[string]interface{}) *UsersLoginInfo {
	newInfo := &UsersLoginInfo{}
	newInfo.Username = (*info)["username"].(string)
	if (*info)["password"] != nil {
		newInfo.Password = (*info)["password"].(string)
	}
	return newInfo
}
func (user *UsersLoginInfo) SetUser() bool {
	var dbError model.DbError

	dbError.MysqlError = SetUserLoginOfMysql(user)
	dbError.RedisError = SetUserOfRedis(user, model.CacheTime)

	return other.ErrorOutput(&dbError)
}
func (user *UsersLoginInfo) UpdateUser() (bool, bool) {
	var dbError model.DbError
	var exist bool
	exist, dbError.MysqlError = UpdateUserLoginOfMysql(user)
	//dbError.RedisError = SetUserOfRedis(user, model.CacheTime)

	return exist, other.ErrorOutput(&dbError)
}
func (user *UsersLoginInfo) DeleteUser() bool {
	var dbError model.DbError
	dbError.MysqlError = DeleteUserLoginOfMysql(user)
	//dbError.RedisError = SetUserOfRedis(&model.UsersInfo{Username: user.Username}, model.GarbageCacheTime)

	return other.ErrorOutput(&dbError)
}
func (user *UsersLoginInfo) GetUser() (bool, bool) {
	var dbError1, dbError2 model.DbError
	var exist bool

	exist, dbError1.RedisError = GetUserOfRedis(user)
	if !exist {
		_, exist, dbError1.MysqlError = GetUserLoginOfMysql(user)
		if exist {
			dbError2.RedisError = SetUserOfRedis(user, model.CacheTime)
			return exist, other.ErrorOutput(&dbError1) && other.ErrorOutput(&dbError2)
		}
		dbError2.RedisError = SetUserOfRedis(&UsersLoginInfo{model.UsersLoginInfo{Username: user.Username}}, model.GarbageCacheTime)
		return exist, other.ErrorOutput(&dbError1) && other.ErrorOutput(&dbError2)
	}
	return exist, other.ErrorOutput(&dbError1) && other.ErrorOutput(&dbError2)
} //获取并更改user信息，返回是否存在和是否error
func LogInUser(username string) bool {
	for i := range 5 {
		err := LogInRedis(username)
		if err == nil {
			return false
		}
		log.Printf("%d times: %v\n", i, err)
		time.Sleep(1 * time.Second)
	}
	return true
}
func LogOutUser(username string) bool {
	for i := range 5 {
		err := LogOutRedis(username)
		if err == nil {
			return false
		}
		log.Printf("%d times: %v\n", i, err)
		time.Sleep(1 * time.Second)
	}
	return true
}
func GetLoginState(username string) (bool, bool) {
	var dbError model.DbError
	var bo bool
	bo, dbError.RedisError = GetLoginStateRedis(username)
	return bo, other.ErrorOutput(&dbError)
} //返回用户登录状态和是否存在error
