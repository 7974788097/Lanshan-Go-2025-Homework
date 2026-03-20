package handler

import (
	"net/http"
	"user-service/kitex_gen/UserService"

	"user-service/kitex_gen/UserService/userservice"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userClient userservice2.Client
}

func NewUserHandler(client userservice2.Client) *UserHandler {
	return &UserHandler{
		userClient: client,
	}
}
func (U *UserHandler) Register(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	resp, err := U.userClient.Register(c.Request.Context(), &userservice.RegisterReq{UserName: req.Username, Password: req.Password})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":     resp.Resp.Code,
		"message":  resp.Resp.Message,
		"httpCode": resp.Resp.HttpCode,
	})
}
func (U *UserHandler) Login(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.ShouldBind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}
	resp, err := U.userClient.Login(c.Request.Context(), &userservice.LoginReq{UserName: req.Username, Password: req.Password})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		c.Abort()
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":     resp.Resp.Code,
		"message":  resp.Resp.Message,
		"httpCode": resp.Resp.HttpCode,
		"Token":    resp.Token,
	})
}
