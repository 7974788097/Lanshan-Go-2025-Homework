package userfollowermanage

import (
	"a_pile_of_shit/app/model"
	"a_pile_of_shit/app/other"
)

var dbContext *model.DBcontext

func ConnectUserFollowerDB(info *model.DBInfo) {
	dbContext = &model.DBcontext{}
	connectUserFollowerMysql(info)
	//connectUserInfoRedis(info)
}

type UsersFollowerInfo struct {
	model.UsersFollowerInfo
	GetWay  string
	GetName []string
}

func Make(info *map[string]interface{}) *UsersFollowerInfo {
	newInfo := &UsersFollowerInfo{}
	if (*info)["follow_username"] != nil {
		newInfo.FollowUsername = (*info)["follow_username"].(string)
	}
	if (*info)["followed_username"] != nil {
		newInfo.FollowedUsername = (*info)["followed_username"].(string)
	}
	if (*info)["get_way"] != nil {
		newInfo.GetWay = (*info)["get_way"].(string)
	}
	return newInfo
}
func (user *UsersFollowerInfo) SetUser() bool {
	var dbError model.DbError

	dbError.MysqlError = SetUserFollowerInfoOfMysql(&(user.UsersFollowerInfo))
	//dbError.RedisError = SetUserOfRedis(user, model.CacheTime)

	return other.ErrorOutput(&dbError)
} //set之前确保关注与被关注的username非空
func (user *UsersFollowerInfo) UpdateUser() (bool, bool) {
	var dbError model.DbError
	var exist bool
	exist, dbError.MysqlError = UpdateUserFollowerInfoOfMysql(&(user.UsersFollowerInfo))
	//dbError.RedisError = SetUserOfRedis(user, model.CacheTime)

	return exist, other.ErrorOutput(&dbError)
}
func (user *UsersFollowerInfo) DeleteUser() bool {
	var dbError model.DbError
	dbError.MysqlError = DeleteUserFollowerInfoOfMysql(&(user.UsersFollowerInfo))
	//dbError.RedisError = SetUserOfRedis(&model.UsersInfo{Username: user.Username}, model.GarbageCacheTime)

	return other.ErrorOutput(&dbError)
}
func (user *UsersFollowerInfo) GetUser() (bool, bool) {
	var dbError model.DbError
	var exist bool
	var end []string
	if user.GetWay == "follow_username" {
		end, exist, dbError.MysqlError = GetUserFollowerInfoOfMysql(&(user.UsersFollowerInfo))
	}
	if user.GetWay == "followed_username" {
		end, exist, dbError.MysqlError = GetUserFollowedInfoOfMysql(&(user.UsersFollowerInfo))
	}
	user.GetName = end
	return exist, other.ErrorOutput(&dbError)
} //获取并更改user信息，返回userinfo,是否存在和是否error
