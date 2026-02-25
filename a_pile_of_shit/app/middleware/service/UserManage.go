package service

import (
	"a_pile_of_shit/app/dao/dbmanage/usermanage"
	"a_pile_of_shit/app/dao/dbmanage/usermanage/userfollowermanage"
	"a_pile_of_shit/app/dao/dbmanage/usermanage/userinfomanage"
	"a_pile_of_shit/app/dao/dbmanage/usermanage/userloginmanage"
	"a_pile_of_shit/app/dao/dbmanage/usermanage/usermutemanage"
	"a_pile_of_shit/app/other"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Register() gin.HandlerFunc {
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
		info["permission"] = 1
		username := info["username"].(string)
		password := info["password"].(string)

		if username == "" || password == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "bad request",
			})
			c.Abort()
			return
		}
		usernameLen := other.CalculateLength(username)
		if usernameLen < 3 {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "username too short",
			})
			c.Abort()
			return
		}
		if usernameLen > 18 {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "username too long",
			})
			c.Abort()
			return
		}
		passwordLen := other.CalculateLength(password)
		if passwordLen < 6 {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "password too short",
			})
		}
		if passwordLen > 18 {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "password too long",
			})
			c.Abort()
			return
		}
		//检验用户是否已经存在
		req := userloginmanage.Make(&info)
		req2 := *req
		var err2 bool
		usermanage.Get(&req2, &bo, &err2)
		if err2 {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Get User Failed",
			})
			c.Abort()
			return
		}
		if bo {
			c.JSON(http.StatusConflict, gin.H{
				"message": "user already exists",
			})
			c.Abort()
			return
		}
		//生成盐值
		salt, err := other.AddSalt()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "salt generation failed",
			})
			c.Abort()
			return
		}
		req.Salt = salt
		//加盐加密
		password = salt + req.Password
		req.Password = other.JiaMi(password)
		//入库
		usermanage.Add(req, &bo)
		var bo2 bool
		var bo3 bool
		info["follow_name"] =
			usermanage.Add(userinfomanage.Make(&info), &bo2)
		usermanage.Add(usermutemanage.Make(&info), &bo3)
		if bo || bo2 || bo3 {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "set user failed",
			})
			c.Abort()
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "register success",
		})
		c.Next()
	}
}
func LogIn() gin.HandlerFunc {
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
		if info["username"] == nil || info["password"] == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "lack info",
			})
			c.Abort()
			return
		}
		username := info["username"].(string)
		password := info["password"].(string)
		if username == "" || password == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "bad request",
			})
			c.Abort()
			return
		}
		req := userloginmanage.Make(&info)
		//检验用户是否存在
		req2 := *req
		inputPassword := req.Password
		var err bool
		usermanage.Get(&req2, &bo, &err)
		if err {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "Get user failed",
			})
			c.Abort()
			return
		}
		if !bo {
			c.JSON(http.StatusConflict, gin.H{
				"message": "user no find",
			})
			c.Abort()
			return
		}
		//检验密码
		if !other.JudgePassword(inputPassword, req2.Password, req2.Salt) {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "wrong password",
			})
			c.Abort()
			return
		}
		c.Set("username", req2.Username)
		c.Set("permission", req2.Permission)
		c.Next()

		//确保token成功生成
		if err, tokenBo := c.Get("tokenError"); tokenBo {
			err = fmt.Errorf("token error: %v", err)
			log.Println(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "token error",
			})
			c.Abort()
			return
		}

		// 返回token
		token, bo := c.Get("token")
		if !bo {
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "login success",
			"token":   token,
		})
	}
}
func Logout() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "logout success",
		})
		c.Next()
	}
}
func ChangePassword() gin.HandlerFunc {
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
		info["username"], _ = c.Get("username")
		if info["old_password"] == "" || info["new_password"] == "" {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "need old password or new password",
			})
			c.Abort()
			return
		}
		OldPassword := info["old_password"].(string)
		NewPassword := info["new_password"].(string)
		//确保有效修改
		if OldPassword == NewPassword {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "The new password is the same as the old password",
			})
			c.Abort()
			return
		}

		var req2 = userloginmanage.Make(&info)
		var exist bool
		usermanage.Get(req2, &exist, &bo)
		if bo {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "get user failed",
			})
			c.Abort()
			return
		}
		if !exist {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "user no find",
			})
			c.Abort()
			return
		}
		if !other.JudgePassword(OldPassword, req2.Password, req2.Salt) {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "password input wrong",
			})
			c.Abort()
			return
		}
		//入库
		req2.Password = other.JiaMi(req2.Salt + NewPassword)
		usermanage.Update(req2, &exist, &bo)
		if bo {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "set user failed",
			})
			c.Abort()
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "change password success",
		})
		c.Next()
	}
}
func ShowUserInfo() gin.HandlerFunc {
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
		var exist bool
		var user = userinfomanage.Make(&info)
		var tempary map[string]interface{}
		tempary = make(map[string]interface{})
		tempary["follow_username"] = info["username"]
		tempary["get_way"] = "followed_username"
		var followInfo = userfollowermanage.Make(&tempary)
		tempary = make(map[string]interface{})
		tempary["followed_username"] = info["username"]
		tempary["get_way"] = "follow_username"
		var followedInfo = userfollowermanage.Make(&tempary)
		usermanage.Get(user, &exist, &bo)
		if bo {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "get user error",
			})
			c.Abort()
			return
		}
		if !exist {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "user no find",
			})
			c.Abort()
			return
		}
		usermanage.Get(followInfo, &exist, &bo)
		if bo {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "get user error",
			})
			c.Abort()
			return
		}
		usermanage.Get(followedInfo, &exist, &bo)
		if bo {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "get user error",
			})
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"followed_names": followInfo.GetName,
			"follow_name":    followedInfo.GetName,
			"permission":     user.Permission,
		})
		c.Next()
	}
} //根据info["get_way"]是"follow_username"还是"followed_name"来查询

func Follow() gin.HandlerFunc {
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
		info["follow_username"], bo = c.Get("username")
		if !bo {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "get info error",
			})
			c.Abort()
			return
		}
		if info["followed_username"] == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "not enter followed username",
			})
			c.Abort()
			return
		}
		user := userfollowermanage.Make(&info)
		usermanage.Add(user, &bo)
		if bo {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "set user error",
			})
			c.Abort()
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "follow success",
		})
		c.Next()
	}
}
func Unfollow() gin.HandlerFunc {
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
		info["follow_username"], bo = c.Get("username")
		if !bo {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "get info error",
			})
			c.Abort()
			return
		}
		if info["followed_username"] == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "not enter followed username",
			})
			c.Abort()
			return
		}
		user := userfollowermanage.Make(&info)
		usermanage.Delete(user, &bo)
		if bo {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "delete user error",
			})
			c.Abort()
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "follow success",
		})
		c.Next()
	}
}
