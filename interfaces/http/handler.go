package http

import (
	"net/http"

	"github.com/DeanThompson/ginpprof"
	"github.com/gin-gonic/gin"

	"mcp-host-demo/infrastructure/common/logit"
	"mcp-host-demo/infrastructure/common/response"
	"mcp-host-demo/interfaces/http/controller"
	"mcp-host-demo/interfaces/middleware"
)

func NewHTTPHandler(
	logger logit.LoggerInterface,
	chat *controller.Chat,
	h *controller.Health,
	response response.HTTPResponseInterface,
) http.Handler {
	r := gin.New()
	r.Use(middleware.Logger(logger, response), middleware.Recovery(logger, true))
	ginpprof.Wrap(r)
	r.GET("/health/liveness", h.Liveness)
	r.GET("/health/readiness", h.Readiness)
	r.POST("/mcp/chat", chat.Chat)
	return r
}
