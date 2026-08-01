package presentation

import (
	"time"

	cachemiddleware "github.com/blocknextai/go-packages/fiber/middleware/cache"
	apiKeysDomain "github.com/blocknextai/platform-api/internal/apikeys/domain/apikeys"
	commonAuth "github.com/blocknextai/platform-api/internal/common/presentation/auth"
	mcpInfrastructure "github.com/blocknextai/platform-api/internal/mcp/infrastructure"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/adaptor"
)

func RegisterPresentation(
	router fiber.Router,
	cacheMiddleware *cachemiddleware.Middleware,
	handlers *mcpInfrastructure.Handlers,
	apiKeyMiddleware *commonAuth.APIKeyMiddleware,
) {
	router.Get(
		"/servers",
		cacheMiddleware.Cache(5*time.Minute),
		NewGetAllServersHandler(handlers.GetAllServers),
	)

	requireAuthentication := apiKeyMiddleware.Authenticate()
	requireMCPScope := apiKeyMiddleware.RequireScope(apiKeysDomain.ScopeMCPInvoke)
	for _, serverHandler := range handlers.ServerHandlers {
		handler := adaptor.HTTPHandler(serverHandler.Handler)
		path := "/" + serverHandler.ID + "/mcp"

		router.Post(path, requireAuthentication, requireMCPScope, handler)
		router.Get(path, requireAuthentication, requireMCPScope, handler)
		router.Delete(path, requireAuthentication, requireMCPScope, handler)
	}
}
