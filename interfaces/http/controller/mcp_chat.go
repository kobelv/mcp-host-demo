package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"mcp-host-demo/application/appservice"
	"mcp-host-demo/application/dto"
	"mcp-host-demo/infrastructure/common/request"
	"mcp-host-demo/infrastructure/common/response"
)

type Chat struct {
	request.BindInterface
	response.HTTPResponseInterface
	app *appservice.ChatAS
}

// NewChat 实例化
func NewChat(req request.BindInterface, res response.HTTPResponseInterface, app *appservice.ChatAS) *Chat {
	return &Chat{BindInterface: req, HTTPResponseInterface: res, app: app}
}

// Chat 入口
func (chat *Chat) Chat(c *gin.Context) {
	dto := dto.ChatDto{}
	if err := chat.Bind(c, &dto); err != nil {
		chat.RenderJSONResponse(c, http.StatusOK, nil, err)
		return
	}
	res, err := chat.app.Chat(c, &dto)
	chat.RenderJSONResponse(c, http.StatusOK, res, err)
}
