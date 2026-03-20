package main

import (
	"client/handler"
	"log"

	"user-service/kitex_gen/UserService/userservice"

	"github.com/cloudwego/kitex/client"
	"github.com/gin-gonic/gin"
	etcd "github.com/kitex-contrib/registry-etcd"
)

func main() {
	r := gin.Default()

	//c, err := userservice.NewClient(
	//	"user-service",
	//	client.WithHostPorts("127.0.0.1:8888"),
	//)
	rResolver, err := etcd.NewEtcdResolver([]string{"127.0.0.1:2379"})
	if err != nil {
		log.Fatal(err)
	}
	userClient, err := userservice2.NewClient(
		"user-service",
		client.WithResolver(rResolver),
	)
	if err != nil {
		log.Fatal(err)
	}
	userHandler := handler.NewUserHandler(userClient)

	api := r.Group("/api")
	{
		userGroup := api.Group("/user")
		{
			userGroup.POST("/register", userHandler.Register)
			userGroup.POST("/login", userHandler.Login)
		}
	}
}
