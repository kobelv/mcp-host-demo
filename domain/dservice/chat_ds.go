package dservice

import (
	"context"
	"mcp-host-demo/domain/adapter"
	"mcp-host-demo/domain/entity"

	"github.com/baidubce/bce-qianfan-sdk/go/qianfan"
	"github.com/mark3labs/mcp-go/mcp"
)

type ChatDS struct {
	chatAdapter      adapter.ChatAdapter
	mcpServerAdapter adapter.McpServerAdapter
}

func NewChatDS(chatAdapter adapter.ChatAdapter, mcpServerAdapter adapter.McpServerAdapter) *ChatDS {
	return &ChatDS{
		chatAdapter:      chatAdapter,
		mcpServerAdapter: mcpServerAdapter,
	}
}

func (ds *ChatDS) Chat(ctx context.Context, entity *entity.ChatInputEntity) (any, error) {
	// #1 new mcp client, list mcp server's tool
	// #2 new LLM, call function call, functions formed by mcp tools
	// #3 mcp client call mcp tool

	tools := ds.mcpServerAdapter.ListMcpTools(ctx)

	fcRes, err := ds.chatAdapter.InvokeFunctionCallArk(ctx, entity.Query, mCPTools2FunctionCall(tools))
	if err != nil {
		return nil, err
	}
	toolRes, err := ds.mcpServerAdapter.InvokeMcpTool(ctx, fcRes)
	if err != nil {
		return nil, err
	}
	return toolRes.Content, nil
}

func mCPTools2FunctionCall(tools *mcp.ListToolsResult) []qianfan.Function {
	if tools == nil || len(tools.Tools) == 0 {
		return nil
	}
	funcs := make([]qianfan.Function, 0)

	for _, v := range tools.Tools {
		funcs = append(funcs, qianfan.Function{
			Name:        v.Name,
			Description: v.Description,
		})
	}
	return funcs
}
