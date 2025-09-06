package adapter

import (
	"context"
	"errors"
	"fmt"
	"mcp-host-demo/domain/adapter"
	"os"
	"time"

	"github.com/baidubce/bce-qianfan-sdk/go/qianfan"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime"
	"github.com/volcengine/volcengine-go-sdk/service/arkruntime/model"

	"mcp-host-demo/infrastructure/common/logit"
)

type ChatAdapterImpl struct {
	logger logit.LoggerInterface
}

const (
	EMMESCOUNT  = 30
	TOKENEXPIRE = 2591990
)

var TokenPrefixPattern = "TokenValue:%s"

func NewChatAdapter(log logit.LoggerInterface) adapter.ChatAdapter {
	return &ChatAdapterImpl{logger: log}
}

func (chat *ChatAdapterImpl) InvokeFunctionCallArk(ctx context.Context, qry string, mcpTools []mcp.Tool) (*qianfan.FunctionCall, error) {
	apiKey := os.Getenv("ARK_API_KEY")
	if apiKey == "" {
		fmt.Println("错误：请设置 ARK_API_KEY 环境变量。")
		return nil, errors.New("-1")
	}

	client := arkruntime.NewClientWithApiKey(
		apiKey,
		arkruntime.WithTimeout(1*time.Minute),
	)

	msg := []*model.ChatCompletionMessage{
		{
			Role: model.ChatMessageRoleUser,
			Content: &model.ChatCompletionMessageContent{
				StringValue: &qry,
			},
		},
	}

	req := model.CreateChatCompletionRequest{
		Model:    "doubao-seed-1-6-flash-250715",
		Messages: msg,
		Tools:    convertToolsToFunctionDefinitions(mcpTools),
		Thinking: &model.Thinking{
			Type: model.ThinkingTypeAuto, // 动态关闭深度思考能力，缩短响应时间
		},
	}

	// 自定义request id
	resp, err := client.CreateChatCompletion(
		ctx,
		req,
		arkruntime.WithCustomHeader(model.ClientRequestHeader, "my-request-id-kobelv"), // 自定义请求头
	)
	if err != nil {
		fmt.Printf("chat error: %v", err)
		return nil, err
	}

	if len(resp.Choices) == 0 {
		fmt.Println("llm returns no choice。")
		return nil, errors.New("-1")
	}

	respMsg := resp.Choices[0].Message

	// 展示模型中间过程的回复内容 (如果存在)
	if respMsg.Content.StringValue != nil && *respMsg.Content.StringValue != "" {
		fmt.Println("llm reply:", *respMsg.Content.StringValue)
		if len(respMsg.ToolCalls) == 0 {
			return nil, errors.New("-1")
		}
	}

	if resp.Choices[0].FinishReason != model.FinishReasonToolCalls || len(respMsg.ToolCalls) == 0 {
		fmt.Println("no function call result:")
		return nil, errors.New("-1")
	}
	if len(respMsg.ToolCalls) > 0 {
		return &qianfan.FunctionCall{
			Name:      respMsg.ToolCalls[0].Function.Name,
			Arguments: respMsg.ToolCalls[0].Function.Arguments,
		}, nil
	}

	return nil, nil
}

func convertToolsToFunctionDefinitions(mcpTools []mcp.Tool) []*model.Tool {
	var arkTools []*model.Tool

	for _, tool := range mcpTools {
		arkTool := &model.Tool{
			Type: "function",
			Function: &model.FunctionDefinition{
				Name:        tool.Name,
				Description: tool.Description,
			},
		}

		arkTool.Function.Parameters = convertParameters(tool.InputSchema)

		arkTools = append(arkTools, arkTool)
	}

	return arkTools
}

func convertParameters(schema mcp.ToolInputSchema) map[string]interface{} {
	params := map[string]interface{}{
		"type":       "object",
		"properties": make(map[string]interface{}),
	}

	if len(schema.Required) > 0 {
		params["required"] = schema.Required
	}
	if len(schema.Properties) > 0 {
		params["properties"] = schema.Properties
	}

	return params
}

func (chat *ChatAdapterImpl) InvokeFunctionCall(ctx context.Context, qry string, functions []qianfan.Function) (*qianfan.FunctionCall, error) {
	qianfan.GetConfig().AccessKey = "tbd"
	qianfan.GetConfig().SecretKey = "tbd"
	llm := qianfan.NewChatCompletion(qianfan.WithModel("ERNIE-4.0-Turbo-8K"))
	msg := []qianfan.ChatCompletionMessage{
		qianfan.ChatCompletionUserMessage(qry),
	}
	mlist := llm.ModelList()
	for _, m := range mlist {
		fmt.Println(m)
	}
	resp, err := llm.Do(context.TODO(),
		&qianfan.ChatCompletionRequest{
			Messages:  msg,
			Functions: functions,
		})
	if err != nil {
		chat.logger.Warn(ctx, "failed to invoke function call: "+err.Error())
		return nil, err
	}

	return resp.FunctionCall, nil
}
