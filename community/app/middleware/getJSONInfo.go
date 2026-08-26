package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		var info map[string]interface{}

		err := c.ShouldBindBodyWithJSON(&info)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "get JSON info error: " + err.Error(),
			})
			c.Abort()
			return
		}
		c.Set("info", info)
		c.Next()
	}
}

//修复中间件解析JSON数据流的问题，而不是err=EOF
