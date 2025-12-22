package api

import (
	"moddle/middleware"
	"moddle/utils"
	"net/http"
	"time"

	"moddle/dao"
	"moddle/model"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var req model.User
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "bad request",
		})
	}
	//确认用户不存在
	if dao.FindUserFromDatabase(req.Username) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "user already exists",
		})
		return
	}
	dao.AddUser(&req)
	c.JSON(http.StatusOK, gin.H{
		"message": "register success",
	})
}

func UpdatePassword(c *gin.Context) {
	type input struct {
		Username         string `json:"username"`
		Newpassword      string `json:"newpassword"`
		Inputoldpassword string `json:"oldpassword"`
	}
	var req input
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "input wrong",
		})
		return
	}
	oldpassword := dao.GetPasswordFromDatabase(req.Username)
	if oldpassword != req.Inputoldpassword { //判断旧密码是否正确
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "old password input wrong",
		})
		return
	}
	dao.UpdatePassword(req.Username, req.Newpassword)
	c.JSON(http.StatusOK, gin.H{
		"message": "update success",
	})
}

func Login(c *gin.Context) {
	var req model.User
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "bad request",
		})
		return
	}
	// 检查用户是否存在且密码是否正确
	if !dao.FindUserFromDatabase(req.Username) {
		dao.AddUserIntoRedis(req.Username, "", 5*time.Second)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "user not found",
		})
		return
	}
	if !dao.JudgePassword(req.Username, req.Password) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "wrong password",
		})
		return
	}

	// 生成jwt token
	token, err := utils.MakeToken(req.Username, time.Now().Add(10*time.Minute))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "internal server error",
		})
		return
	}
	// 返回token
	c.JSON(http.StatusOK, gin.H{
		"message": "login",
		"token":   token,
	})
}
func InitRouter_gin() {
	r := gin.Default()
	r.GET("/ping", middleware.JudgeToken(), middleware.Example2(), ping1)
	r.POST("/login", Login)
	r.POST("/register", Register)
	r.POST("/updatepassword", middleware.JudgeToken(), UpdatePassword)
	r.POST("/delete")

	_ = r.Run(":8080")
}
