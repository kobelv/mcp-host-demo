package adapter

import (
	"context"
	"errors"
	"fmt"
	"mcp-host-demo/domain/adapter"
	"os"
	"time"

	"github.com/baidubce/bce-qianfan-sdk/go/qianfan"
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

func (chat *ChatAdapterImpl) InvokeFunctionCallArk(ctx context.Context, qry string, functions []qianfan.Function) (*qianfan.FunctionCall, error) {
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
	tools := make([]*model.Tool, 0)
	for _, v := range functions {
		tools = append(tools, &model.Tool{
			Type: model.ToolTypeFunction,
			Function: &model.FunctionDefinition{
				Name:        v.Name,
				Description: v.Description,
			},
		})
	}

	req := model.CreateChatCompletionRequest{
		Model:    "doubao-seed-1-6-flash-250715", // 替换为Model ID，请从文档获取 https://www.volcengine.com/docs/82379/1330310
		Messages: msg,
		Tools:    tools,
		Thinking: &model.Thinking{
			Type: model.ThinkingTypeAuto, // 关闭深度思考能力
		},
	}

	// 自定义request id
	resp, err := client.CreateChatCompletion(
		ctx,
		req,
		arkruntime.WithCustomHeader(model.ClientRequestHeader, "my-request-id-a1b2c3"), // 自定义请求头
	)
	if err != nil {
		fmt.Printf("standard chat error: %v", err)
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
		return nil, errors.New("-1")
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

func (chat *ChatAdapterImpl) InvokeFunctionCall(ctx context.Context, qry string, functions []qianfan.Function) (*qianfan.FunctionCall, error) {
	qianfan.GetConfig().AccessKey = "PwohSMMMbgjO2CLQ4jMYy1MY"
	qianfan.GetConfig().SecretKey = "K6FM61SkrXQ2AGa8niSTG9sr535qbTNn"
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
		chat.logger.Warn(ctx, "failed to call function call: "+err.Error())
		return nil, err
	}

	return resp.FunctionCall, nil
}
