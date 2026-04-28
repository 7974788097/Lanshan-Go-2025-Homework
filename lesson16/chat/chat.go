package chat

import (
	"context"
	"errors"
	"io"
	"time"

	"github.com/cloudwego/eino-ext/components/model/openai"
	"github.com/cloudwego/eino/components/prompt"
	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent/react"
	"github.com/cloudwego/eino/schema"
)

func CreateAiModel(ctx context.Context, ModelName string, BaseURL string, APIKey string) (*openai.ChatModel, error) {
	return openai.NewChatModel(ctx, &openai.ChatModelConfig{
		Model:   ModelName,
		BaseURL: BaseURL,
		APIKey:  APIKey,
	})
}
func CreateAiAgent(ctx context.Context, ModelName string, BaseURL string, APIKey string, toolInfo []tool.BaseTool, maxStep int) (*react.Agent, error) {
	aiModel, err := CreateAiModel(ctx, ModelName, BaseURL, APIKey)
	if err != nil {
		return nil, err
	}
	return react.NewAgent(ctx, &react.AgentConfig{
		ToolCallingModel: aiModel,
		ToolsConfig: compose.ToolsNodeConfig{
			Tools: toolInfo,
		},
		MaxStep: maxStep,
	})
}
func CreateHistory(maxChatTurns int) []*schema.Message {
	return make([]*schema.Message, 0, (maxChatTurns+1)*2)
}
func CreateFormate(role string, needTodo string) *prompt.DefaultChatTemplate {
	systemMessage := "你是一个" + role + ",你需要" + needTodo + ",你使用中文进行交流，绝对不能夹杂其他语言"
	return prompt.FromMessages(schema.FString, schema.SystemMessage(systemMessage), schema.UserMessage("{usermessage}"), schema.MessagesPlaceholder("chat_history", true))
}
func CreateMessage(formate *prompt.DefaultChatTemplate, userMessage string, history []*schema.Message) []*schema.Message {
	message, err := formate.Format(context.Background(), map[string]any{"usermessage": userMessage, "chat_history": history})
	if err != nil {
		panic(err)
	}
	return message
}
func AddHistory(history []*schema.Message, userMessage string, assistantMessage *schema.Message) []*schema.Message {
	return append(history, schema.UserMessage(userMessage), assistantMessage)
}
func CleanHistory(keepNumber int, history []*schema.Message) []*schema.Message {
	length := len(history)
	if length <= 2*keepNumber {
		return history
	}
	history = history[length-2*keepNumber:]
	return history
}
func AgentStreamChat(ctx context.Context, message []*schema.Message, Agent *react.Agent, outTime time.Duration, recept chan string) (*schema.Message, error) {
	var output string
	var reader *schema.StreamReader[*schema.Message]
	var err error
	timeCtx, cancel := context.WithCancel(ctx)
	resetChan := make(chan struct{}, 1)
	errChan := make(chan error, 1)
	resetTimeout := func() {
		select {
		case resetChan <- struct{}{}:
		default:
		}
	}
	timer := time.NewTimer(outTime)
	go func() {
		for {
			select {
			//case <-time.After(outTime):
			//	cancel()
			//	return
			case <-timer.C:
				cancel()
				return
			case <-timeCtx.Done():
				return
			case <-resetChan:
			}
			timer.Reset(outTime)
		}
	}()
	go func() {
		for output == "" {

			select {
			case <-timeCtx.Done():
				return
			default:
			}
			reader, err = Agent.Stream(timeCtx, message)
			if err != nil {
				reader.Close()
				errChan <- err
				return
			}
			for {
				select {
				case <-timeCtx.Done():
					reader.Close()
					errChan <- errors.New("timeout")
					return
				default:
				}
				msg, err := reader.Recv()
				if err == io.EOF {
					break
				}
				if err != nil {
					reader.Close()
					errChan <- err
					return
				}
				if msg.Content != "" {
					output += msg.Content
					recept <- msg.Content
					resetTimeout()
				}
			}
			reader.Close()
		}
		cancel()
	}()
	select {
	case err = <-errChan:
		cancel()
		close(resetChan)
		close(errChan)
		close(recept)
		return nil, err
	case <-timeCtx.Done():
	}
	close(resetChan)
	close(errChan)
	close(recept)
	return schema.AssistantMessage(output, nil), nil
}
