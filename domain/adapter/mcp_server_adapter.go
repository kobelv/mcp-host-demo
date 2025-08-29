package adapter

import (
	"context"

	"github.com/baidubce/bce-qianfan-sdk/go/qianfan"
	"github.com/mark3labs/mcp-go/mcp"
)

type McpServerAdapter interface {
	ListMcpTools(ctx context.Context) *mcp.ListToolsResult
	InvokeMcpTool(ctx context.Context, fcRes *qianfan.FunctionCall) (*mcp.CallToolResult, error)
}
