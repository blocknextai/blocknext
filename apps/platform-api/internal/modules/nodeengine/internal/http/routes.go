package http

import (
	"time"

	"github.com/gofiber/fiber/v3"

	cachemiddleware "github.com/blocknextai/go-packages/fiber/middleware/cache"
	nodeEngineHTTPCredentials "github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/http/credentials"
	nodeEngineHTTPNodes "github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/http/nodes"
	nodeEngineHTTPWebhooks "github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/http/webhooks"
	nodeEngineUseCases "github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/usecases"
)

func RegisterRoutes(
	router fiber.Router,
	cacheMiddleware *cachemiddleware.Middleware,
	useCases *nodeEngineUseCases.Services,
) {
	nodeEngineRouterGroup := router.Group("/node-engine")

	nodeEngineRouterGroup.Get(
		"/nodes",
		cacheMiddleware.Cache(5*time.Minute),
		nodeEngineHTTPNodes.NewGetAllNodesHandler(useCases.Nodes),
	)

	nodeEngineRouterGroup.Get(
		"/credentials",
		cacheMiddleware.Cache(5*time.Minute),
		nodeEngineHTTPCredentials.NewGetAllCredentialsHandler(useCases.Credentials),
	)

	nodeEngineRouterGroup.Get(
		"/credentials/:id",
		cacheMiddleware.Cache(5*time.Minute),
		nodeEngineHTTPCredentials.NewGetCredentialByIDHandler(useCases.Credentials),
	)

	nodeEngineRouterGroup.Get(
		"/webhook-sources",
		cacheMiddleware.Cache(5*time.Minute),
		nodeEngineHTTPWebhooks.NewGetAllWebhookSourcesHandler(useCases.Adapters),
	)

	nodeEngineRouterGroup.Get(
		"/trigger-variables",
		cacheMiddleware.Cache(5*time.Minute),
		nodeEngineHTTPWebhooks.NewGetAllTriggerVariablesHandler(useCases.Adapters),
	)
}
