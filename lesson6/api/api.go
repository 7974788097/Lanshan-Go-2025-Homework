package api

import (
	"moddle/middleware"
	"net/http"
	"time"

	"moddle/dao"
	"moddle/model"
	"moddle/utils"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var req model.User
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "bad request",
		})
	}
	// 如果用户存在，这里这种是用户名可以一致的，即只要密码不一致就视为不同用户
	if dao.FindUser(req.Username, req.Password) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "user already exists",
		})
		return
	}
	err := dao.AddUser(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Database Update error",
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "register success",
	})
}

func UpdatePassword(username string, newpassword string, inputoldpassword string, c *gin.Context) {
	oldpassword := dao.SelectPasswordFromUsername(username)
	if oldpassword != inputoldpassword { //判断旧密码是否正确
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "bad request",
		})
		return
	}
	err := dao.UpdatePassword(username, newpassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"messgae": "Database Update error",
		})
		return
	}
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
	if !dao.FindUser(req.Username, req.Password) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "user not found",
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
	_ = r.Run(":8080")
}
