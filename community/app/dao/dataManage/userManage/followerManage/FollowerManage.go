package followermanage

import "Homework_LanShan/new_shit/app/model"

type FollowerStruct struct {
	model.UsersFollowerInfo
	GetInfo []FollowerStruct
}

func NewFollowerStruct(info *map[string]any) *FollowerStruct {
	newStruct := new(FollowerStruct)
	if (*info)["follow_username"] != nil {
		newStruct.FollowUsername = (*info)["follow_username"].(string)
	}
	if (*info)["followed_username"] != nil {
		newStruct.FollowedUsername = (*info)["followed_username"].(string)
	}
	return newStruct
}
func (U *FollowerStruct) Add(context *model.DBcontext) error {
	err := U.addInMysql(context.Mysql)
	return err
}
