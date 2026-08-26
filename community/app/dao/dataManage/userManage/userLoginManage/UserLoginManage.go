package userloginmmanage

import "Homework_LanShan/new_shit/app/model"

type UserLoginInfoStruct struct {
	model.UsersLoginInfo
}

func NewUserLoginInfoStruct(info *map[string]any) *UserLoginInfoStruct {
	newStruct := new(UserLoginInfoStruct)
	if (*info)["username"] != nil {
		newStruct.Username = (*info)["username"].(string)
	}
	if (*info)["password"] != nil {
		newStruct.Password = (*info)["password"].(string)
	}
	if (*info)["salt"] != nil {
		newStruct.Salt = (*info)["salt"].(string)
	}
	if (*info)["permission"] != nil {
		newStruct.Permission = (*info)["permission"].(uint8)
	}
	return newStruct
}
