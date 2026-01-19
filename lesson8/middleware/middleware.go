package middleware

import (
	"errors"
	"fmt"
	"moddle/model"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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
		parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return model.SecretOfToken, nil
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
		c.Set("username", claims["username"])
		c.Next()
	}
}

//func JudgeToken() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		token := c.GetHeader("Token")
//		if token == "" {
//			c.JSON(http.StatusUnauthorized, gin.H{
//				"message": "no find Token,please login",
//			})
//			c.Abort()
//			return
//		}
//		parts := strings.Split(token, "#")
//		if len(parts) != 2 {
//			c.JSON(http.StatusUnauthorized, gin.H{
//				"message": "Token's format wrong",
//			})
//			c.Abort()
//		}
//		username := parts[0]
//		expireStr := parts[1]
//		expireTime, err := strconv.ParseInt(expireStr, 10, 64)
//		if err != nil || time.Now().Unix() > expireTime {
//			c.JSON(http.StatusUnauthorized, gin.H{
//				"message": "Token invalid or expired",
//			})
//			c.Abort()
//			return
//		}
//		c.Set("username", username)
//		c.Next()
//	}
//}
