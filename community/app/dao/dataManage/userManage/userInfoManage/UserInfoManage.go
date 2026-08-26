package userinfomanage

import "Homework_LanShan/new_shit/app/model"

type UserInfoStruct struct {
	model.UsersInfo
}

func NewUserInfoStruct(info *map[string]any) *UserInfoStruct {
	newStruct := new(UserInfoStruct)
	if (*info)["username"] != nil {
		newStruct.Username = (*info)["username"].(string)
	}
	if (*info)["permission"] != nil {
		newStruct.Permission = (*info)["permission"].(uint8)
	}
	return newStruct
}
