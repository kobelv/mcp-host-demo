package dservice

import (
	"context"
	"mcp-host-demo/domain/adapter"
	"mcp-host-demo/domain/entity"
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
	// #1 list mcp server's tools
	// #2 new LLM, invoke function call, compose functions by mcp server‘s tools via step 1
	// #3 mcp client invoke mcp server tool call via step 2

	tools := ds.mcpServerAdapter.ListMcpTools(ctx)

	fcRes, err := ds.chatAdapter.InvokeFunctionCallArk(ctx, entity.Query, tools.Tools)
	if err != nil {
		return nil, err
	}
	toolRes, err := ds.mcpServerAdapter.InvokeMcpTool(ctx, fcRes)
	if err != nil {
		return nil, err
	}
	return toolRes.Content, nil
}
