package main

import (
	"log"
	"user-service/dao"
	handler2 "user-service/handler"
	userservice "user-service/kitex_gen/UserService/userservice"
)

func main() {
	dbInfo := dao.GetDBInfo()

	dbContext := dao.ConnectUserDB(dbInfo)
	dao.AutoUserManage(dbContext.Mysql)
	log.Println("success")
	handler := &handler2.UserServiceImpl{
		DBContext: dbContext,
	}

	// ========== 3. 创建并启动 Kitex 服务 ==========
	svr := userservice.NewServer(handler)
	//svr := userservice.NewServer(&UserServiceImpl{dbContext})

	err := svr.Run()

	if err != nil {
		log.Println(err.Error())
	}
}
