package api

import (
	"a_pile_of_shit/app/dao/dbmanage"
	"a_pile_of_shit/app/middleware"
	"a_pile_of_shit/app/middleware/service"
	"a_pile_of_shit/app/other"

	"github.com/gin-gonic/gin"
)

func Begin() {
	middleware.GetTokenPassword()
	other.MakeSnowNode()
	dao.ConnectDatabase()
}
func InitRouter() {
	limiter := middleware.NewIpRateLimiter(2, 6)
	g := gin.Default()
	g.Use(middleware.Limiter(limiter), middleware.GetInfo())
	needLogin := g.Group("")
	needLogin.Use(middleware.AnalysisToken(), middleware.JudgeLoginState())

	comment := needLogin.Group("/comment")
	{
		comment.POST("/set", service.JudgeMuteState(), service.SetComment())
		comment.GET("/get", service.GetComment())
		comment.POST("/delete", service.DeleteComment())
	}
	message := needLogin.Group("/message")
	{
		message.POST("/set", service.JudgeMuteState(), service.SetMessage())
		message.POST("/send", service.SendMessage())
		message.POST("/update", service.UpdateMessage())
		message.POST("/delete", service.DeleteMessage())
		message.GET("/search", service.SearchMessage())
		message.GET("/get", service.GetMessage(), middleware.JudgePermission())
		message.POST("/setPermission", middleware.SetMessagePermission())
	}
	user := needLogin.Group("/user")
	{
		g.POST("/user/register", service.Register())
		g.POST("/user/login", service.LogIn(), middleware.LogInUser(), middleware.MakeToken())
		user.POST("/logout", middleware.LogOutUser(), service.Logout())
		user.POST("/changePassword", service.ChangePassword())
		user.GET("/get", service.ShowUserInfo())
		user.POST("/setPermission", middleware.SetUserPermission(), middleware.LogOutUser())
		user.POST("/follow", service.Follow())
		user.POST("/unFollow", service.Unfollow())

		user.POST("/setMute", middleware.JudgeManagePermission(), service.SetMuteLabel())
		user.POST("/releaseMute", middleware.JudgeManagePermission(), service.ReleaseMuteState())
	}
	err := g.Run(":8080")
	if err != nil {
		panic(err)
	}
}
