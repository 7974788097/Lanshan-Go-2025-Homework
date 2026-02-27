package service

import (
	"a_pile_of_shit/app/dao/dbmanage/datamanage"
	"a_pile_of_shit/app/dao/dbmanage/datamanage/commentmanage"
	"a_pile_of_shit/app/other"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetComment() gin.HandlerFunc {
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
		if info["parent_node_id"] == nil || info["message_id"] == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "not get the parent node id",
			})
			c.Abort()
			return
		}
		//输入检验
		if info["parent_node_id"] == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "not get the parent node id",
			})
			c.Abort()
			return
		}
		if info["comment_content"] == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "not get the content",
			})
			c.Abort()
			return
		}
		//长度限制
		commentContentLen := other.CalculateLength(info["comment_content"].(string))
		if commentContentLen > 10000 {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "comment length too large",
			})
			c.Abort()
			return
		}
		//生成实例
		var path string
		info["comment_id"] = other.MakeID()*10 + 1
		Comment := commentmanage.Make(&info)
		var exist bool
		//为父节点增加信息
		//父节点为评论时
		if int64(info["parent_node_id"].(float64))%10 == 1 {
			var a map[string]interface{}
			a["comment_id"] = info["parent_node_id"]
			a["get_way"] = "comment_id"
			parentComment := commentmanage.Make(&a)
			datamanage.Get(parentComment, &exist, &bo)
			if bo {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "get comment error",
				})
				c.Abort()
				return
			}
			if !exist {
				c.JSON(http.StatusBadRequest, gin.H{
					"message": "not find the parent node id",
				})
				c.Abort()
				return
			}
			path = parentComment.Path
			parentComment.ChildNodeID = append(parentComment.ChildNodeID, Comment.CommentID)
			datamanage.Update(parentComment, &exist, &bo)
			if bo {
				c.JSON(http.StatusInternalServerError, gin.H{
					"message": "update parent comment error",
				})
				c.Abort()
				return
			}
			if !exist {
				c.JSON(http.StatusBadRequest, gin.H{
					"message": "not find the parent node id",
				})
				c.Abort()
				return
			}
		} else if int64(info["parent_node_id"].(float64))%10 == 0 { //父节点为文章时
			path = fmt.Sprintf("%x", int64(info["parent_node_id"].(float64))/10) + "\\"
		} else {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "parent node id error",
			})
			c.Abort()
			return
		}
		//处理评论信息
		Comment.Path = fmt.Sprintf("%s%x\\", path, Comment.CommentID)
		datamanage.Add(Comment, &bo)
		if bo {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "set comment error",
			})
			c.Abort()
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "set comment success",
		})
		c.Next()
	}
}
func DeleteComment() gin.HandlerFunc {
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
		if info["comment_path"] == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "not get the comment path",
			})
			c.Abort()
			return
		}
		//生成实例
		comment := commentmanage.Make(&info)
		datamanage.Delete(comment, &bo)
		if bo {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "delete comment error",
			})
			c.Abort()
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "delete comment success",
		})
		c.Next()
	}
}
func GetComment() gin.HandlerFunc {
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
		if info["get_way"] == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "lack get_way",
			})
			c.Abort()
			return
		}
		if (info["comment_id"] == nil || info["message_id"] == nil) && info["get_way"].(string) != "username" {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "lack search id",
			})
			c.Abort()
			return
		}
		//文章id查询时
		if info["message_id"] != nil && info["get_way"].(string) == "message_id" {
			info["comment_path"] = fmt.Sprintf("%x\\", int64(info["message_id"].(float64)))
			info["get_way"] = "path"
		}
		getWay := info["get_way"].(string)
		if getWay != "comment_id" && getWay != "path" && getWay != "username" {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "get_way error,only support comment_id ,message_id ,path and username",
			})
		}
		//获取评论信息
		Comment := commentmanage.Make(&info)
		var exist bool
		datamanage.Get(Comment, &exist, &bo)
		if bo {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "get comment error",
			})
			c.Abort()
			return
		}
		if !exist {
			c.JSON(http.StatusOK, gin.H{
				"message": "comment not exist",
			})
			c.Next()
			return
		}
		//返回数据
		c.JSON(http.StatusOK, gin.H{
			"message":  "get comment success",
			"comments": Comment.Info,
		})
		c.Next()
	}
} //可以根据message_id,username,comment_id进行查询
