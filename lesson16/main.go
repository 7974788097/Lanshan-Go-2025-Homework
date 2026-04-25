package main

import (
	"aiagent/api"
	"aiagent/chat"
	"aiagent/config"
	"aiagent/tool"
	"context"
	"time"
)

func main() {
	info := config.GetInfo()
	tools := tool.InitTools()
	aiAgent, err := chat.CreateAiAgent(context.Background(), info.ModelName, info.BaseURL, info.APIKey, tools, info.MaxThinkStep)
	if err != nil {
		panic(err)
	}
	formate := chat.CreateFormate(info.AiRole, info.NeedTodo)
	apiModel := api.NewModelInfo(aiAgent, formate, 10*time.Second, info.MaxChatTurns)
	api.InitRouter(apiModel)
}
