package service

import (
	"a_pile_of_shit/app/dao/dbmanage/datamanage"
	"a_pile_of_shit/app/dao/dbmanage/datamanage/messagemanage"
	"a_pile_of_shit/app/other"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetMessage() gin.HandlerFunc {
	return func(c *gin.Context) {
		rawInfo, bo := c.Get("info")
		if !bo {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Get info error",
			})
			c.Abort()
			return
		}
		info := rawInfo.(map[string]interface{})
		info["username"], bo = c.Get("username")
		if !bo {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Get info error",
			})
			c.Abort()
			return
		}
		//输入检验
		if info["message_content"] == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "lack message content",
			})
			c.Abort()
			return
		}
		if info["message_name"] == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "lack message name",
			})
			c.Abort()
			return
		}
		//长度限制
		messageContentLen := other.CalculateLength(info["message_content"].(string))
		if messageContentLen > 50000 {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Message content length too large",
			})
			c.Abort()
			return
		}
		messageNameLen := other.CalculateLength(info["message_name"].(string))
		if messageNameLen > 30 {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Message name length too large",
			})
			c.Abort()
			return
		}
		//生成实例
		message := messagemanage.Make(&info)
		message.MessageID = other.MakeID() * 10
		message.Permission = 0
		message.Status = 1
		datamanage.Add(message, &bo)
		if bo {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Set message error",
			})
			c.Abort()
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "Set message success",
		})
		c.Next()
	}
}
func GetMessage() gin.HandlerFunc {
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
		info["username"], bo = c.Get("username")
		if !bo {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Get info error",
			})
			c.Abort()
			return
		}
		//输入检验
		if info["get_way"] == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "lack get_way",
			})
			c.Abort()
			return
		}
		if info["get_way"].(string) != "message_id" && info["get_way"].(string) != "username" {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "get_way error",
			})
			c.Abort()
			return
		}
		if info["message_id"] == nil && info["get_way"].(string) == "message_id" {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "lack message id",
			})
			c.Abort()
			return
		}
		//生成实例
		message := messagemanage.Make(&info)
		var exist bool
		datamanage.Get(message, &exist, &bo)
		if bo {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "get message error",
			})
			c.Abort()
			return
		}
		if !exist {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "message id not exist",
			})
			c.Abort()
			return
		}
		c.Set("needPermission", message.Permission)
		c.Set("writer", message.Username)
		c.Next()
		JudgePermissionStatus, bo := c.Get("JudgePermissionStatus")
		if !bo {
			if !exist {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "get permission status error",
				})
				c.Abort()
				return
			}
		}
		switch JudgePermissionStatus.(int) {
		case 0:
			c.JSON(http.StatusForbidden, gin.H{
				"message": "permission not enough",
			})
		case 1:
			c.JSON(http.StatusOK, gin.H{
				"message": message.Info,
			})
		case 2:
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "get message error",
			})
		}
	}
} //通过message_id或username获取文章信息
func UpdateMessage() gin.HandlerFunc {
	return func(c *gin.Context) {
		rawInfo, bo := c.Get("info")
		if !bo {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Get info error",
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
		if info["message_content"] == nil && info["message_name"] == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "lack useful info",
			})
			c.Abort()
			return
		}
		message := messagemanage.Make(&info)
		var exist bool
		if info["message_content"] != nil {
			message.Content = info["message_content"].(string)
		}
		if info["message_name"] != nil {
			message.MessageName = info["message_name"].(string)
		}
		datamanage.Update(message, &exist, &bo)
		if bo {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Update message error",
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
			"message": "Update message success",
		})
		c.Next()
	}
}
func DeleteMessage() gin.HandlerFunc {
	return func(c *gin.Context) {
		rawInfo, bo := c.Get("info")
		if !bo {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Get info error",
			})
			c.Abort()
			return
		}
		info := rawInfo.(map[string]interface{})
		//输入检验
		if info["message_id"] == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "lack message id",
			})
			c.Abort()
			return
		}
		//创建实例
		var a map[string]interface{}
		a = make(map[string]interface{})
		a["get_way"] = "path"
		a["comment_path"] = fmt.Sprintf("%x\\", int64(info["message_id"].(float64)))
		comment := messagemanage.Make(&a)
		datamanage.Delete(comment, &bo)
		if bo {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Delete comment error",
			})
			c.Abort()
			return
		}
		message := messagemanage.Make(&info)
		datamanage.Delete(message, &bo)
		if bo {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Delete message error",
			})
			c.Abort()
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "Delete message success",
		})
		c.Next()
	}
}
func SearchMessage() gin.HandlerFunc {
	return func(c *gin.Context) {
		rawInfo, bo := c.Get("info")
		if !bo {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Get info error",
			})
			c.Abort()
			return
		}
		info := rawInfo.(map[string]interface{})
		info["username"] = info["writer"]

		//输入检验
		if info["get_way"] == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "lack get_way",
			})
			c.Abort()
			return
		}
		if info["message_id"] == nil && info["writer"] == nil && info["message_name"] == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "lack search info",
			})
			c.Abort()
			return
		}
		message := messagemanage.Make(&info)
		var exist bool
		datamanage.Get(message, &exist, &bo)
		if bo {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Get message error",
			})
			c.Abort()
			return
		}
		if !exist {
			c.JSON(http.StatusOK, gin.H{
				"message": "not find message",
			})
			c.Next()
			return
		}
		lenth := len(message.Info)
		messageId := make([]int64, 0, lenth)
		for _, i := range message.Info {
			if i.Status == 2 {
				messageId = append(messageId, i.MessageID)
			}
		}
		c.JSON(http.StatusOK, gin.H{
			"message": gin.H{
				"message_id": messageId,
			},
		})
		c.Next()
	}
} //能根据message_id,writer或者message_name查找已发送的文章
func SendMessage() gin.HandlerFunc {
	return func(c *gin.Context) {
		rawInfo, bo := c.Get("info")
		if !bo {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Get info error",
			})
			c.Abort()
			return
		}
		info := rawInfo.(map[string]interface{})
		info["username"], bo = c.Get("username")
		if !bo {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Get info error",
			})
			c.Abort()
			return
		}
		info["get_way"] = "message_id"
		//输入检验
		if info["message_id"] == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "lack message id",
			})
			c.Abort()
			return
		}
		//创建实例
		message := messagemanage.Make(&info)
		var exist bool
		datamanage.Get(message, &exist, &bo)
		if bo {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Get message error",
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
		if message.Username != info["username"].(string) {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "you are not allowed to send message",
				"reason":  "you are not the writer",
			})
			c.Abort()
			return
		}
		message.Status = 2
		datamanage.Update(message, &exist, &bo)
		if bo {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Update message error",
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
			"message": "sent message success",
		})
		c.Next()
	}
}
