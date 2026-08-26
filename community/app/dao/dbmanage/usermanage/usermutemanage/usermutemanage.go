package usermutemanage

import (
	"a_pile_of_shit/app/model"
	"a_pile_of_shit/app/other"
)

var dbContext *model.DBcontext

func ConnectUserMuteDB(info *model.DBInfo) {
	dbContext = &model.DBcontext{}
	connectUserMuteMysql(info)
	//connectUserInfoRedis(info)
}

type UsersMuteInfo struct {
	model.UsersMuteInfo
}

func Make(info *map[string]interface{}) *UsersMuteInfo {
	newInfo := &UsersMuteInfo{}
	newInfo.Username = (*info)["username"].(string)
	if (*info)["mute_lable"] != nil {
		newInfo.MuteLabel = (*info)["mute_label"].([2]int)
	}
	if (*info)["mute_reason"] != nil {
		newInfo.MuteReason = (*info)["mute_reason"].(string)
	}
	return newInfo
}
func (user *UsersMuteInfo) SetUser() bool {
	var dbError model.DbError

	dbError.MysqlError = SetUserMuteInfoOfMysql(user)
	//dbError.RedisError = SetUserOfRedis(user, model.CacheTime)

	return other.ErrorOutput(&dbError)
}
func (user *UsersMuteInfo) UpdateUser() (bool, bool) {
	var dbError model.DbError
	var exist bool
	exist, dbError.MysqlError = UpdateUserMuteInfoOfMysql(user)
	//dbError.RedisError = SetUserOfRedis(user, model.CacheTime)

	return exist, other.ErrorOutput(&dbError)
}
func (user *UsersMuteInfo) DeleteUser() bool {
	var dbError model.DbError
	dbError.MysqlError = DeleteUserMuteInfoOfMysql(user)
	//dbError.RedisError = SetUserOfRedis(&model.UsersInfo{Username: user.Username}, model.GarbageCacheTime)

	return other.ErrorOutput(&dbError)
}
func (user *UsersMuteInfo) GetUser() (bool, bool) {
	var dbError model.DbError
	var exist bool
	_, exist, dbError.MysqlError = GetUserMuteInfoOfMysql(user)
	return exist, other.ErrorOutput(&dbError)
} //获取并更改user信息，返回userinfo,是否存在和是否error
