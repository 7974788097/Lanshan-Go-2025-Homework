package api

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ping1(c *gin.Context) {
	name, exist := c.Get("username")
	if !exist {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "username is nil",
		})
		return
	}
	information := fmt.Sprintf("pong %s", name)
	c.JSON(http.StatusOK, gin.H{
		"message": information,
	})
}
