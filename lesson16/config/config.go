package config

import (
	"aiagent/other"
	"strconv"
)

type ConfigInfo struct {
	ModelName    string
	BaseURL      string
	APIKey       string
	AiRole       string
	NeedTodo     string
	MaxThinkStep int
	MaxChatTurns int
}

func GetInfo() ConfigInfo {
	var raw int64 //临时变量
	var err error
	newConfigInfo := ConfigInfo{}
	newConfigInfo.ModelName = other.GetEnv("MODEL_NAME")
	newConfigInfo.BaseURL = other.GetEnv("BASE_URL")
	newConfigInfo.APIKey = other.GetEnv("API_KEY")
	newConfigInfo.AiRole = other.GetEnv("AI_ROLE")
	newConfigInfo.NeedTodo = other.GetEnv("NEED_TODO")
	raw, err = strconv.ParseInt(other.GetEnv("MAX_THINK_STEP"), 10, 64)
	if err != nil {
		panic(".env中MAX_THINK_STEP配置错误：" + err.Error())
	}
	newConfigInfo.MaxThinkStep = int(raw)
	raw, err = strconv.ParseInt(other.GetEnv("MAX_CHAT_TURNS"), 10, 64)
	if err != nil {
		panic(".env中MAX_CHAT_TURNS配置错误：" + err.Error())
	}
	newConfigInfo.MaxChatTurns = int(raw)
	return newConfigInfo
}
