//go:build wireinject
// +build wireinject

package interfaces

import (
	"context"

	"github.com/google/wire"

	"mcp-host-demo/application/appservice"
	"mcp-host-demo/application/dto"
	"mcp-host-demo/domain/dservice"
	"mcp-host-demo/infrastructure/adapter"
	"mcp-host-demo/infrastructure/common/cache"
	"mcp-host-demo/infrastructure/common/db"
	"mcp-host-demo/infrastructure/common/logit"
	"mcp-host-demo/infrastructure/common/request"
	"mcp-host-demo/infrastructure/common/response"
	"mcp-host-demo/interfaces/http"
	"mcp-host-demo/interfaces/http/controller"
)

func NewApp(ctx context.Context) (*app, error) {
	panic(wire.Build(wire.NewSet(
		loadAppConf,
		logit.NewServiceLoggerConf,
		logit.NewServiceLogger,

		db.NewDB,
		cache.NewRedis,
		response.NewHTTPResponseWriter,
		request.NewRequest,

		controller.NewChat,
		controller.NewHealth,
		http.NewHTTPHandler,
		newHTTPServer,

		appservice.NewChatAS,

		dto.NewChatDtoConv,
		dservice.NewChatDS,
		adapter.NewChatAdapter,
		adapter.NewMcpServerAdapter,

		wire.Struct(new(app), "*"),
	)))
}
