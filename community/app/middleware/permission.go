package middleware

import (
	"a_pile_of_shit/app/dao/dbmanage/datamanage"
	"a_pile_of_shit/app/dao/dbmanage/datamanage/messagemanage"
	"a_pile_of_shit/app/dao/dbmanage/usermanage"
	"a_pile_of_shit/app/dao/dbmanage/usermanage/userinfomanage"
	"a_pile_of_shit/app/dao/dbmanage/usermanage/userloginmanage"
	"net/http"

	"github.com/gin-gonic/gin"
)

//普通权限：1，会员权限：2，管理员权限：3

func JudgePermission() gin.HandlerFunc {
	return func(c *gin.Context) {
		needPermission, bo := c.Get("needPermission")
		permission, bo2 := c.Get("permission")
		writer, bo3 := c.Get("writer")
		username := c.GetString("username")
		if !(bo && bo2 && bo3) {
			c.Set("JudgePermissionStatus", 2)
			c.Abort()
			return
		}
		if needPermission.(uint8) > permission.(uint8) && writer.(string) != username {
			c.Set("JudgePermissionStatus", 0)
			c.Abort()
			return
		}
		c.Set("JudgePermissionStatus", 1)
		c.Next()
	}
}
func SetMessagePermission() gin.HandlerFunc {
	return func(c *gin.Context) {
		rawInfo, bo := c.Get("info")
		if !bo {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "get info error",
			})
			c.Abort()
			return
		}
		info := rawInfo.(map[string]interface{})
		info["get_way"] = "message_id"
		//输入检验
		if info["message_id"] == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "lack message id",
			})
			c.Abort()
			return
		}
		if info["message_permission"] == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "lack message permission",
			})
			c.Abort()
			return
		}
		message := messagemanage.Make(&info)
		var exist bool
		message.Permission = info["message_permission"].(uint8)
		datamanage.Update(message, &exist, &bo)
		if bo {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "update permission error",
			})
			c.Abort()
			return
		}
		if !exist {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "not find message",
			})
			c.Abort()
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "change permission success",
		})
		c.Next()
	}
}
func SetUserPermission() gin.HandlerFunc {
	return func(c *gin.Context) {
		rawInfo, bo := c.Get("info")
		if !bo {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "get info error",
			})
			c.Abort()
			return
		}
		info := rawInfo.(map[string]interface{})
		//输入检验
		if info["goal_username"] == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "lack goal username",
			})
			c.Abort()
			return
		}
		if info["goal_permission"] == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "lack goal permission",
			})
			c.Abort()
			return
		}

		chacheinfo := make(map[string]interface{})
		chacheinfo["username"] = info["goal_username"].(string)
		chacheinfo["permission"] = uint8(info["goal_permission"].(float64))
		user := userinfomanage.Make(&chacheinfo)
		userlogin := userloginmanage.Make(&chacheinfo)
		var exist bool
		usermanage.Update(user, &exist, &bo)
		if bo {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "update user error",
			})
			c.Abort()
			return
		}
		if !exist {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "not find user",
			})
			c.Abort()
			return
		}
		usermanage.Update(userlogin, &exist, &bo)
		if bo {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "update user error",
			})
			c.Abort()
			return
		}
		if !exist {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "not find user",
			})
			c.Abort()
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "change permission success",
		})
		c.Next()
	}
}
func JudgeManagePermission() gin.HandlerFunc {
	return func(c *gin.Context) {
		permission, bo := c.Get("permission")
		if !bo {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "get permission error",
			})
			c.Abort()
			return
		}
		if permission.(uint8) < 3 {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "you are not manager",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
