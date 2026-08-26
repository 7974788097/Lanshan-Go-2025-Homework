package usermanage

import (
	"a_pile_of_shit/app/dao/dbmanage/usermanage/userfollowermanage"
	"a_pile_of_shit/app/dao/dbmanage/usermanage/userinfomanage"
	"a_pile_of_shit/app/dao/dbmanage/usermanage/userloginmanage"
	"a_pile_of_shit/app/dao/dbmanage/usermanage/usermutemanage"
	"a_pile_of_shit/app/model"
)

func ConnectUserDB(info *model.DBInfo) {
	userinfomanage.ConnectUserInfoDB(info)
	userloginmanage.ConnectUserLoginDB(info)
	usermutemanage.ConnectUserMuteDB(info)
	userfollowermanage.ConnectUserFollowerDB(info)
}

type UserOperation interface {
	SetUser() bool
	UpdateUser() (bool, bool)
	GetUser() (bool, bool)
	DeleteUser() bool
}

func Add(user UserOperation, bo *bool) UserOperation {
	*bo = user.SetUser()
	return user
}
func Delete(user UserOperation, bo *bool) UserOperation {
	*bo = user.DeleteUser()
	return user
}
func Get(user UserOperation, exist *bool, bo *bool) UserOperation {
	*exist, *bo = user.GetUser()
	return user
}
func Update(user UserOperation, exist *bool, bo *bool) UserOperation {
	*exist, *bo = user.UpdateUser()
	return user
}
func LogInUser(username string) bool {
	return userloginmanage.LogInUser(username)
}
func LogOutUser(username string) bool {
	return userloginmanage.LogOutUser(username)
}
func GetLoginState(username string) (bool, bool) {
	return userloginmanage.GetLoginState(username)
}
