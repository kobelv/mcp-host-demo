package adapter

import (
	"context"

	"github.com/baidubce/bce-qianfan-sdk/go/qianfan"
	"github.com/mark3labs/mcp-go/mcp"
)

type ChatAdapter interface {
	InvokeFunctionCall(ctx context.Context, qry string, functions []qianfan.Function) (*qianfan.FunctionCall, error)
	InvokeFunctionCallArk(ctx context.Context, qry string, functions []mcp.Tool) (*qianfan.FunctionCall, error)
}
