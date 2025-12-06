package middleware

import (
	"moddle/dao"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func JudgeToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Token")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "no find Token,please login",
			})
			c.Abort()
			return
		}
		parts := strings.Split(token, "#")
		if len(parts) != 2 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Token's format wrong",
			})
		}
		username := parts[0]
		expireStr := parts[1]
		expireTime, err := strconv.ParseInt(expireStr, 10, 64)
		if err != nil || time.Now().Unix() > expireTime {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Token invalid or expired",
			})
			c.Abort()
			return
		}
		if _, exist := dao.Database[username]; !exist {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Username not exist",
			})
			c.Abort()
			return
		}
		c.Set("username", username)
		c.Next()
	}
}
func Example2() gin.HandlerFunc {
	return func(c *gin.Context) {

	}
}
