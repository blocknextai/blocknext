package presentation

import (
	"time"

	cachemiddleware "github.com/blocknextai/go-packages/fiber/middleware/cache"
	"github.com/blocknextai/go-packages/rbac"
	apiKeysInfrastructure "github.com/blocknextai/platform-api/internal/apikeys/infrastructure"
	apikeys "github.com/blocknextai/platform-api/internal/apikeys/presentation/apikeys"
	commonPresentationAuth "github.com/blocknextai/platform-api/internal/common/presentation/auth"
	"github.com/gofiber/fiber/v3"
)

func RegisterOrganizationAPIKeysPresentation(
	router fiber.Router,
	authMiddleware *commonPresentationAuth.AuthMiddleware,
	handlers *apiKeysInfrastructure.Handlers,
) {
	apiKeysRouterGroup := router.Group("/organizations/:organizationId/api-keys")

	apiKeysRouterGroup.Get(
		"/",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.ReadAPIKeyPermission),
		apikeys.NewGetAllOrganizationAPIKeysHandler(handlers.GetAllAPIKeys),
	)

	apiKeysRouterGroup.Post(
		"/",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.CreateAPIKeyPermission),
		apikeys.NewCreateOrganizationAPIKeyHandler(handlers.CreateAPIKey),
	)

	apiKeysRouterGroup.Post(
		"/:apiKeyId/regenerate",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.UpdateAPIKeyPermission),
		apikeys.NewRegenerateOrganizationAPIKeyHandler(handlers.RegenerateAPIKey),
	)

	apiKeysRouterGroup.Delete(
		"/:apiKeyId",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.DeleteAPIKeyPermission),
		apikeys.NewDeleteOrganizationAPIKeyHandler(handlers.DeleteAPIKey),
	)
}

func RegisterAPIKeysPresentation(
	router fiber.Router,
	cacheMiddleware *cachemiddleware.Middleware,
	handlers *apiKeysInfrastructure.Handlers,
) {
	apiKeysRouterGroup := router.Group("/api-keys")

	apiKeysRouterGroup.Get(
		"/scopes",
		cacheMiddleware.Cache(1*time.Hour),
		apikeys.NewGetAPIKeyScopesHandler(handlers.GetScopes),
	)
}

func RegisterPresentation(
	router fiber.Router,
	authMiddleware *commonPresentationAuth.AuthMiddleware,
	cacheMiddleware *cachemiddleware.Middleware,
	handlers *apiKeysInfrastructure.Handlers,
) {
	RegisterOrganizationAPIKeysPresentation(router, authMiddleware, handlers)
	RegisterAPIKeysPresentation(router, cacheMiddleware, handlers)
}
