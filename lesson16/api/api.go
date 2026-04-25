package api

import (
	"aiagent/chat"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/flow/agent/react"
	"github.com/cloudwego/eino/schema"
	"github.com/gin-gonic/gin"
)

type ModelInfo struct {
	AiAgent        *react.Agent
	MessageFormate *prompt.DefaultChatTemplate
	UserInfo       sync.Map
	OutTime        time.Duration
	MaxChatTurns   int
}

func NewModelInfo(Agent *react.Agent, Formate *prompt.DefaultChatTemplate, OutTime time.Duration, MaxChatTurns int) *ModelInfo {
	return &ModelInfo{
		AiAgent:        Agent,
		MessageFormate: Formate,
		OutTime:        OutTime,
		MaxChatTurns:   MaxChatTurns,
	}
}
func InitRouter(info *ModelInfo) {
	gin.SetMode(gin.DebugMode)
	g := gin.Default()

	g.POST("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	g.Use(func() gin.HandlerFunc {
		return func(c *gin.Context) {
			//这里面用来获取用户信息，省略
		}
	}())
	g.POST("/chat", info.StreamChat)

	if err := g.Run(":8080"); err != nil {
		panic(err)
	}
}

func (i *ModelInfo) StreamChat(c *gin.Context) {
	var err error
	UserID := c.GetString("user_id")
	input := c.GetString("message")
	if input == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"code":    114514,
			"message": "lack message",
		})
		return
	}
	var history []*schema.Message
	rawHistory, exit := i.UserInfo.Load(UserID)
	if !exit {
		history = chat.CreateHistory(i.MaxChatTurns)
	} else {
		history = rawHistory.([]*schema.Message)
	}
	message := chat.CreateMessage(i.MessageFormate, input, history)
	recept := make(chan string, 15)
	var output *schema.Message
	go func() {
		output, err = chat.AgentStreamChat(c.Request.Context(), message, i.AiAgent, i.OutTime, recept)
	}()

	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	for content := range recept {
		_, _ = c.Writer.WriteString(fmt.Sprintf("data: %s\n\n", content))
		c.Writer.Flush()
	}
	if err != nil {
		_, _ = c.Writer.WriteString(fmt.Sprintf("data: %s\n\n", "\n发生错误了"))
		c.Writer.Flush()
		output = schema.AssistantMessage("发生错误了", nil)
	}
	history = chat.AddHistory(history, input, output)
	history = chat.CleanHistory(i.MaxChatTurns, history)
	i.UserInfo.Store(UserID, history)
}
