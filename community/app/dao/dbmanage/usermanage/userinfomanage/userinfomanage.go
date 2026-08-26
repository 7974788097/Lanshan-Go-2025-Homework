package userinfomanage

import (
	"a_pile_of_shit/app/model"
	"a_pile_of_shit/app/other"
)

var dbContext *model.DBcontext

func ConnectUserInfoDB(info *model.DBInfo) {
	dbContext = &model.DBcontext{}
	connectUserInfoMysql(info)
	//connectUserInfoRedis(info)
}

type UsersInfo struct {
	model.UsersInfo
}

func Make(info *map[string]interface{}) *UsersInfo {
	newInfo := &UsersInfo{}
	newInfo.Username = (*info)["username"].(string)
	if (*info)["user_permission"] != nil {
		newInfo.Permission = (*info)["user_permission"].(uint8)
	}
	return newInfo
}
func (user *UsersInfo) SetUser() bool {
	var dbError model.DbError

	dbError.MysqlError = SetUserInfoOfMysql(user)
	//dbError.RedisError = SetUserOfRedis(user, model.CacheTime)

	return other.ErrorOutput(&dbError)
}
func (user *UsersInfo) UpdateUser() (bool, bool) {
	var dbError model.DbError
	var exist bool
	exist, dbError.MysqlError = UpdateUserInfoOfMysql(user)
	//dbError.RedisError = SetUserOfRedis(user, model.CacheTime)

	return exist, other.ErrorOutput(&dbError)
}
func (user *UsersInfo) DeleteUser() bool {
	var dbError model.DbError
	dbError.MysqlError = DeleteUserInfoOfMysql(user)
	//dbError.RedisError = SetUserOfRedis(&model.UsersInfo{Username: user.Username}, model.GarbageCacheTime)

	return other.ErrorOutput(&dbError)
}
func (user *UsersInfo) GetUser() (bool, bool) {
	var dbError model.DbError
	var exist bool
	_, exist, dbError.MysqlError = GetUserInfoOfMysql(user)
	return exist, other.ErrorOutput(&dbError)
	//var dbError1, dbError2 model.DbError
	//var bo bool
	//
	//bo, dbError1.RedisError = GetUserOfRedis(user)
	//if !bo {
	//	_, bo, dbError1.MysqlError = GetUserOfMysql(user)
	//	if bo {
	//		dbError2.RedisError = SetUserOfRedis(user, model.CacheTime)
	//		return *user, bo, other.ErrorOutput(&dbError1) && other.ErrorOutput(&dbError2)
	//	}
	//	dbError2.RedisError = SetUserOfRedis(&model.UsersLoginInfo{Username: user.Username}, model.GarbageCacheTime)
	//	return UsersInfo{}, bo, other.ErrorOutput(&dbError1) && other.ErrorOutput(&dbError2)
	//}
	//return *user, bo, other.ErrorOutput(&dbError1) && other.ErrorOutput(&dbError2)
} //获取并更改user信息，返回userinfo,是否存在和是否error
