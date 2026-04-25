package tool

import (
	"context"
	"time"

	"github.com/bytedance/sonic"
	agenttool "github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
)

func InitTools() []agenttool.BaseTool {
	return []agenttool.BaseTool{
		&GetTime{},
	}
}

type GetTime struct{}

func (g *GetTime) Info(_ context.Context) (*schema.ToolInfo, error) {
	return &schema.ToolInfo{
		Name: "get_time",
		Desc: "【当用户询问当前时间的时候必须调用该工具】获取当先的系统时间",
	}, nil
}
func (g *GetTime) InvokableRun(_ context.Context, _ string, _ ...agenttool.Option) (string, error) {
	type Resp struct {
		Time string `json:"time"`
	}
	resp := &Resp{
		Time: time.Now().Format("2006-01-02 15:04:05"),
	}
	return sonic.MarshalString(resp)
}
