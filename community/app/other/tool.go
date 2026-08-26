package other

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"time"
	"unicode/utf8"

	"github.com/redis/go-redis/v9"
)

func AddSalt() (string, error) {
	saltBytes := make([]byte, 32)
	_, err := rand.Read(saltBytes)
	if err != nil {
		log.Println("AddSalt:", err)
		return "", fmt.Errorf("salt generation failed: %v", err)
	}
	return hex.EncodeToString(saltBytes), nil
} //生成盐值，未加密

func JiaMi(input string) string {
	h := sha256.New()
	h.Write([]byte(input))
	return hex.EncodeToString(h.Sum(nil))
} //返回加密后的字符串

func JudgePassword(input string, password string, salt string) bool {
	return hmac.Equal([]byte(password), []byte(JiaMi(salt+input)))
} //返回bool判断输入的原始密码是否和加盐加密的用户密码匹配

func RedisHSetEX(redisctx *redis.Client, ctx context.Context, cacheTime time.Duration, field string, info map[string]interface{}) error {
	tx := redisctx.TxPipeline()
	for key, value := range info {
		tx.HSet(ctx, field, key, value)
	}
	tx.Expire(ctx, field, cacheTime)
	_, err := tx.Exec(ctx)
	return err
} //使用事务，HSet同时设置缓存时间
func CalculateLength[T []any | string](input T) int {
	switch v := any(input).(type) {
	case []any:
		return len(v)
	case string:
		return utf8.RuneCountInString(v)
	}
	return 0
}
