package middleware

import (
	"a_pile_of_shit/app/dao/dbmanage/usermanage"
	"net/http"

	"github.com/gin-gonic/gin"
)

func LogInUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		rawinfo, bo := c.Get("info")
		if !bo {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "get info error",
			})
			c.Abort()
			return
		}
		info := rawinfo.(map[string]interface{})

		username := info["username"].(string)
		err := usermanage.LogInUser(username)
		if err {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Login user error ",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
func LogOutUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		rawinfo, bo := c.Get("info")
		if !bo {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "get info error",
			})
			c.Abort()
			return
		}
		info := rawinfo.(map[string]interface{})

		username := info["username"].(string)
		usermanage.LogOutUser(username)
		c.Next()
	}
}
func JudgeLoginState() gin.HandlerFunc {
	return func(c *gin.Context) {
		rawusername, bo := c.Get("username")
		if !bo {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "get info error",
			})
			c.Abort()
			return
		}
		username := rawusername.(string)
		exist, bo := usermanage.GetLoginState(username)
		if bo {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "login state get error",
			})
			c.Abort()
			return
		}
		if !exist {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "User not Login",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
