package utils

import (
	"strconv"
	"strings"
	"time"
)

func MakeToken(username string, expireTime time.Time) (string, error) {
	// 核心逻辑：用特殊分隔符（#）拼接 用户名 + 过期时间戳（秒级）
	// 分隔符选一个业务中不会出现的字符，避免解析时出错
	expireTimestamp := strconv.FormatInt(expireTime.Unix(), 10)
	token := strings.Join([]string{username, expireTimestamp}, "#")
	return token, nil
}
