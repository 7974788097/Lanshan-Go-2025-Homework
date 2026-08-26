package service

import (
	"a_pile_of_shit/app/dao/dbmanage/usermanage"
	"a_pile_of_shit/app/dao/dbmanage/usermanage/usermutemanage"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func JudgeMuteState() gin.HandlerFunc {
	return func(c *gin.Context) {
		rawinfo, bo := c.Get("info")
		if !bo {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Get info error",
			})
			c.Abort()
			return
		}
		info := rawinfo.(map[string]interface{})
		info["username"], bo = c.Get("username")
		if !bo {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Get info error",
			})
			c.Abort()
			return
		}
		user := usermutemanage.Make(&info)
		if user.MuteLabel[0] == 2 || user.MuteLabel[1] >= int(time.Now().Unix()) {
			c.JSON(http.StatusForbidden, gin.H{
				"message":    "being muted",
				"endTime":    time.Unix(int64(user.MuteLabel[1]), 0).Format("2006-01-02 15:04:05"),
				"muteReason": user.MuteReason,
			})
			c.Abort()
			return
		}
		c.Next()
	}
}
func SetMuteLabel() gin.HandlerFunc {
	return func(c *gin.Context) {
		rawinfo, bo := c.Get("info")
		if !bo {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Get info error",
			})
			c.Abort()
			return
		}
		info := rawinfo.(map[string]interface{})
		//输入检查
		if info["mute_reason"] == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "lack mute reason",
			})
			c.Abort()
			return
		}
		if info["mute_time"] == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "lack mute time",
			})
			c.Abort()
		}
		if info["mute_username"] == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "lack mute username",
			})
			c.Abort()
			return
		}
		var exist bool
		var a map[string]interface{}
		a["username"] = info["mute_username"]
		muteUser := usermutemanage.Make(&a)
		muteUser.MuteLabel[0] = 2
		muteUser.MuteLabel[1] = int(time.Now().Add(time.Duration(int64(info["mute_time"].(float64))) * time.Second).Unix())
		muteUser.MuteReason = info["mute_reason"].(string)
		usermanage.Update(muteUser, &exist, &bo)
		if bo {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "set user failed",
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
			"message": "mute user success",
		})
		c.Next()
	}
}
func ReleaseMuteState() gin.HandlerFunc {
	return func(c *gin.Context) {
		rawinfo, bo := c.Get("info")
		if !bo {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Get info error",
			})
			c.Abort()
			return
		}
		info := rawinfo.(map[string]interface{})
		permission, bo := c.Get("permission")
		if !bo {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Get info error",
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
		//输入检验
		if info["mute_username"] == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "lack mute username",
			})
			c.Abort()
			return
		}
		var a map[string]interface{}
		a["username"] = info["mute_username"]
		muteUser := usermutemanage.Make(&a)
		var exist bool

		muteUser.MuteLabel[0] = 1
		muteUser.MuteLabel[1] = 0
		muteUser.MuteReason = "无"
		usermanage.Update(muteUser, &exist, &bo)
		if bo {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "set user failed",
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
			"message": "mute user success",
		})
		c.Next()
	}
}
