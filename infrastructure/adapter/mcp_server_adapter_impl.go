package adapter

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"mcp-host-demo/domain/adapter"

	"mcp-host-demo/infrastructure/common/logit"

	"github.com/baidubce/bce-qianfan-sdk/go/qianfan"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/client/transport"
	"github.com/mark3labs/mcp-go/mcp"
)

type McpServerAdapterImpl struct {
	logger     logit.LoggerInterface
	c          *client.Client
	serverInfo *mcp.InitializeResult
}

func NewMcpServerAdapter(ctx context.Context, log logit.LoggerInterface) adapter.McpServerAdapter {
	sse, err := transport.NewSSE("http://10.12.214.141:8088/sse")
	if err != nil {
		fmt.Printf("failed to create SSE %v\n", err)
	}

	if err := sse.Start(ctx); err != nil {
		fmt.Printf("failed to start SSE %v\n", err)
	}

	c := client.NewClient(sse)
	// defer c.Close()

	req := mcp.InitializeRequest{}
	req.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	req.Params.ClientInfo = mcp.Implementation{
		Name:    "I'm MCP Client",
		Version: "1.0.0",
	}
	req.Params.Capabilities = mcp.ClientCapabilities{}
	serverInfo, err := c.Initialize(ctx, req)
	if err != nil {
		fmt.Printf("failed to init client %v\n", err)
	}
	fmt.Printf("connected to server: %s, version is %s\n", serverInfo.ServerInfo.Name, serverInfo.ServerInfo.Version)
	if err := c.Ping(ctx); err != nil {
		fmt.Printf("failed to ping %v\n", err)
	}

	return &McpServerAdapterImpl{
		logger:     log,
		c:          c,
		serverInfo: serverInfo,
	}
}

func (adapter *McpServerAdapterImpl) ListMcpTools(ctx context.Context) *mcp.ListToolsResult {
	if adapter.serverInfo.Capabilities.Tools != nil {
	}

	toolsReq := mcp.ListToolsRequest{}
	toolRes, err := adapter.c.ListTools(ctx, toolsReq)
	if err != nil {
		adapter.logger.Info(ctx, "failed to list tools")
	}

	return toolRes
}

func (adapter *McpServerAdapterImpl) InvokeMcpTool(ctx context.Context, fcRes *qianfan.FunctionCall) (*mcp.CallToolResult, error) {
	if fcRes == nil {
		return nil, nil
	}
	if fcRes.Name != "greet" {
		return nil, errors.New("this tool is not supported yet")
	}

	arg := make(map[string]any)
	json.Unmarshal([]byte(fcRes.Arguments), &arg)

	callToolRequest := mcp.CallToolRequest{
		Params: struct {
			Name      string    `json:"name"`
			Arguments any       `json:"arguments,omitempty"`
			Meta      *mcp.Meta `json:"_meta,omitempty"`
		}{
			Name: fcRes.Name,
			Arguments: map[string]any{
				"name": arg["name"].(string),
			},
		},
	}

	result, err := adapter.c.CallTool(ctx, callToolRequest)
	if err != nil {
		fmt.Printf("failed to invoke tool %v\n", err)
		return nil, err
	}

	return result, nil
}
