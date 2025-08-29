package dto

import (
	"mcp-host-demo/domain/entity"
)

type ChatDtoConv struct{}

func NewChatDtoConv() *ChatDtoConv {
	return &ChatDtoConv{}
}

func (conv *ChatDtoConv) Dto2Entity(qry *ChatDto) *entity.ChatInputEntity {
	return &entity.ChatInputEntity{
		Query:   qry.Query,
		Remarks: qry.Remarks,
	}
}
