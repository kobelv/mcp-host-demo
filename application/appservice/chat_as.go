package appservice

import (
	"context"
	"mcp-host-demo/application/dto"
	"mcp-host-demo/domain/dservice"
)

type ChatAS struct {
	ChatDS  *dservice.ChatDS
	DtoConv *dto.ChatDtoConv
}

func NewChatAS(mcpDS *dservice.ChatDS) *ChatAS {
	return &ChatAS{
		ChatDS: mcpDS,
	}
}

func (as *ChatAS) Chat(ctx context.Context, qry *dto.ChatDto) (any, error) {
	return as.ChatDS.Chat(ctx, as.DtoConv.Dto2Entity(qry))
}
