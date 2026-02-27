package other

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"reflect"
	"strconv"
	"time"
	"unicode/utf8"

	"github.com/bwmarrin/snowflake"
	"github.com/redis/go-redis/v9"
)

func GetenvDefault(key string, def string) string {
	a := os.Getenv(key)
	if a == "" {
		a = def
	}
	return a
} //获取env文件的值，未获取便设置默认值
var globalSnowNode *snowflake.Node

func MakeSnowNode() {
	nodeString := GetenvDefault("SNOW_NODE_ID", "1")
	node, err := strconv.Atoi(nodeString)
	if err != nil {
		panic(fmt.Errorf("make snow node id error : %v", err))
	}
	globalSnowNode, err = snowflake.NewNode(int64(node))
	if err != nil {
		panic(fmt.Errorf("make snow node error : %v", err))
	}
	beginningTimeString := GetenvDefault("BEGINNING_TIME", "0")
	beginningTime, err := strconv.Atoi(beginningTimeString)
	if err != nil {
		panic(fmt.Errorf("get beginning time error : %v", err))
	}
	snowflake.Epoch = int64(beginningTime)
}
func MakeID() int64 {
	return globalSnowNode.Generate().Int64()
} //生成唯一的16-19位ID
func ErrorOutput(DbErr interface{}) bool {
	var bo = false
	val := reflect.ValueOf(DbErr).Elem()
	field := val.NumField()
	for i := range field {
		fieldval := val.Field(i)
		err, ok := fieldval.Interface().(error)
		if ok && err != nil {
			fieldval2 := err.Error()
			if fieldval2 != "" {
				bo = true
				log.Println(fieldval)
			}
		}

	}
	return bo
} //统一打印数据库操作的错误日志，无error返回false

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

func RedisHSetEX(redisctx *redis.Client, ctx context.Context, cacheTime time.Duration, username string, info map[string]interface{}) error {
	tx := redisctx.TxPipeline()
	for key, value := range info {
		tx.HSet(ctx, username, key, value)
	}
	tx.Expire(ctx, username, cacheTime)
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
