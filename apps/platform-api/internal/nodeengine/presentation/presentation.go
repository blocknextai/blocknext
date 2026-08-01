package presentation

import (
	"time"

	cachemiddleware "github.com/blocknextai/go-packages/fiber/middleware/cache"
	nodeEngineInfrastructure "github.com/blocknextai/platform-api/internal/nodeengine/infrastructure"
	nodeEnginePresentationCredentials "github.com/blocknextai/platform-api/internal/nodeengine/presentation/credentials"
	nodeEnginePresentationNodes "github.com/blocknextai/platform-api/internal/nodeengine/presentation/nodes"
	nodeEnginePresentationWebhooks "github.com/blocknextai/platform-api/internal/nodeengine/presentation/webhooks"
	"github.com/gofiber/fiber/v3"
)

func RegisterPresentation(
	router fiber.Router,
	cacheMiddleware *cachemiddleware.Middleware,
	handlers *nodeEngineInfrastructure.Handlers,
) {
	nodeEngineRouterGroup := router.Group("/node-engine")

	nodeEngineRouterGroup.Get(
		"/nodes",
		cacheMiddleware.Cache(5*time.Minute),
		nodeEnginePresentationNodes.NewGetAllNodesHandler(handlers.GetAllNodes),
	)

	nodeEngineRouterGroup.Get(
		"/credentials",
		cacheMiddleware.Cache(5*time.Minute),
		nodeEnginePresentationCredentials.NewGetAllCredentialsHandler(handlers.GetAllCredentials),
	)

	nodeEngineRouterGroup.Get(
		"/credentials/:id",
		cacheMiddleware.Cache(5*time.Minute),
		nodeEnginePresentationCredentials.NewGetCredentialByIDHandler(handlers.GetCredentialByID),
	)

	nodeEngineRouterGroup.Get(
		"/webhook-sources",
		cacheMiddleware.Cache(5*time.Minute),
		nodeEnginePresentationWebhooks.NewGetAllWebhookSourcesHandler(handlers.GetAllWebhookSources),
	)

	nodeEngineRouterGroup.Get(
		"/trigger-variables",
		cacheMiddleware.Cache(5*time.Minute),
		nodeEnginePresentationWebhooks.NewGetAllTriggerVariablesHandler(handlers.GetAllTriggerVariables),
	)
}
