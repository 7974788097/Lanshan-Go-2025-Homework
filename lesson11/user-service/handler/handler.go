package handler

import (
	"context"
	userservice "user-service/kitex_gen/UserService"
	"user-service/model"
	"user-service/service"
)

// UserServiceImpl implements the last service interface defined in the IDL.
type UserServiceImpl struct {
	DBContext *model.DBcontext
}

// Register implements the UserServiceImpl interface.
func (s *UserServiceImpl) Register(ctx context.Context, req *userservice.RegisterReq) (resp *userservice.RegisterResp, err error) {
	var info = map[string]interface{}{"username": req.UserName, "password": req.Password}
	//省略一大堆数据检验
	user := service.NewUserInfoService(s.DBContext, &info)
	return user.Register()
}

// Login implements the UserServiceImpl interface.
func (s *UserServiceImpl) Login(ctx context.Context, req *userservice.LoginReq) (resp *userservice.LoginResp, err error) {
	var info = map[string]interface{}{"username": req.UserName, "password": req.Password}
	//省略一大堆数据检验
	user := service.NewUserInfoService(s.DBContext, &info)
	return user.Login()
}
