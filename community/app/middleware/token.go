package middleware

import (
	"a_pile_of_shit/app/dao/dbmanage/usermanage"
	"a_pile_of_shit/app/dao/dbmanage/usermanage/userloginmanage"
	"a_pile_of_shit/app/model"
	"a_pile_of_shit/app/other"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

var tokenPassword string

func GetTokenPassword() {
	tokenPassword = other.GetenvDefault("TOKEN_PASSWORD", "123456")
}
func MakeToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.MustGet("username")
		var info map[string]interface{}
		info = make(map[string]interface{})
		info["username"] = username
		user := userloginmanage.Make(&info)
		var bo, exist bool
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
				"message": "user not exist",
			})
			c.Abort()
			return
		}
		claims := jwt.MapClaims{
			"username":   username,
			"permission": user.Permission,
			"expiretime": time.Now().Add(model.CacheTime).Unix(),
		}
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString([]byte(tokenPassword))
		if err != nil {
			c.Set("tokenError", err)
			c.Abort()
			return
		}
		c.Set("token", tokenString)
		c.Next()
	}
} //生成包含username，permission和expire time的token
func AnalysisToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Token")
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "no find Token,please login",
			})
			c.Abort()
			return
		}
		parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(tokenPassword), nil
		})
		if err != nil {
			var errMsg string
			// 精准区分JWT库定义的标准错误类型
			switch {
			// 1. Token已过期（最常见，优先提示）
			case errors.Is(err, jwt.ErrTokenExpired):
				errMsg = "invalid token: expired"
			// 2. Token签名错误（被篡改/秘钥错误）
			case errors.Is(err, jwt.ErrSignatureInvalid):
				errMsg = "invalid token: signature error"
			// 3. Token尚未生效（nbf字段设置了生效时间）
			case errors.Is(err, jwt.ErrTokenNotValidYet):
				errMsg = "invalid token: not valid yet"
			// 4. Token格式错误（不是3段式/缺少.分隔符）
			case strings.Contains(err.Error(), "malformed"):
				errMsg = "invalid token: malformed format"
			// 5. 算法不匹配（你自定义的算法错误）
			case strings.Contains(err.Error(), "unexpected signing method"):
				errMsg = "invalid token: unsupported algorithm"
			// 6. 其他未知错误（保留原标识，便于定位）
			default:
				errMsg = "invalid token(1): " + err.Error()
			}
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": errMsg,
			})
			c.Abort()
			return
		}
		claims, ok := parsedToken.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "invalid token:claims parse failed",
			})
			c.Abort()
			return
		}
		expFloat, ok2 := claims["expiretime"].(float64)
		if !ok2 {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "invalid token: Missing exp",
			})
			c.Abort()
			return
		}
		exp := int64(expFloat)
		if time.Now().Unix() > exp {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "invalid token: expired",
			})
			c.Abort()
			return
		}
		if !parsedToken.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "invalid token(2)",
			})
			c.Abort()
			return
		}
		if claims["username"] == nil || claims["permission"] == nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "invalid token(3)",
			})
			c.Abort()
			return
		}
		c.Set("username", claims["username"])
		permission := claims["permission"].(float64)
		c.Set("permission", uint8(permission))
		c.Next()
	}
} //解析token是否合法以及是否与前端username参数相同，将username和permission存入context
