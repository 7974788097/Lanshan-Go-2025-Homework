package service

import (
	"user-service/dao"
	"user-service/dao/userLoginManage"
	common "user-service/kitex_gen/Common"
	userservice "user-service/kitex_gen/UserService"
	"user-service/model"
)

type UserInfoService struct {
	UserDBcontext *model.DBcontext
	Info          *map[string]interface{}
}

func NewUserInfoService(dbContext *model.DBcontext, Info *map[string]any) *UserInfoService {
	UserInfo := &UserInfoService{}
	UserInfo.UserDBcontext = dbContext
	UserInfo.Info = Info
	return UserInfo
}
func (U *UserInfoService) Register() (*userservice.RegisterResp, error) {
	var err error
	userLogin := userLoginManage.Make(U.Info)
	err = dao.Add(userLogin, U.UserDBcontext)
	if err != nil {
		return &userservice.RegisterResp{
			Resp: &common.Resp{
				2001,
				500,
				"register fail",
			},
		}, err
	}
	return &userservice.RegisterResp{
		Resp: &common.Resp{
			0,
			200,
			"register ok",
		},
	}, nil
}
func (U *UserInfoService) Login() (*userservice.LoginResp, error) {
	var err error
	var bo bool
	userLogin := userLoginManage.Make(U.Info)
	bo, err = dao.Get(userLogin, U.UserDBcontext)
	if err != nil {
		return &userservice.LoginResp{
			Resp: &common.Resp{
				2001,
				500,
				"login fail",
			},
		}, err
	}
	if !bo {
		return &userservice.LoginResp{
			Resp: &common.Resp{
				1002,
				400,
				"user not exist",
			},
		}, nil
	}
	if (*U.Info)["password"].(string) != userLogin.Password {
		return &userservice.LoginResp{
			Resp: &common.Resp{
				1003,
				400,
				"password fail",
			},
		}, nil
	}
	Token := "this is just a Token"
	return &userservice.LoginResp{
		Resp: &common.Resp{
			0,
			200,
			"success",
		},
		Token: Token,
	}, nil
}
