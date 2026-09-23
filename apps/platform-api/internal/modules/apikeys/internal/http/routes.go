package http

import (
	"time"

	"github.com/gofiber/fiber/v3"

	cachemiddleware "github.com/blocknextai/go-packages/fiber/middleware/cache"
	"github.com/blocknextai/go-packages/rbac"
	commonAuth "github.com/blocknextai/platform-api/internal/common/auth"
	apikeys "github.com/blocknextai/platform-api/internal/modules/apikeys/internal/http/apikeys"
	apiKeysUseCases "github.com/blocknextai/platform-api/internal/modules/apikeys/internal/usecases"
)

func RegisterOrganizationAPIKeysRoutes(
	router fiber.Router,
	authMiddleware *commonAuth.AuthMiddleware,
	useCases *apiKeysUseCases.Services,
) {
	apiKeysRouterGroup := router.Group("/organizations/:organizationId/api-keys")

	apiKeysRouterGroup.Get(
		"/",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.ReadAPIKeyPermission),
		apikeys.NewGetAllOrganizationAPIKeysHandler(useCases.APIKeys),
	)

	apiKeysRouterGroup.Post(
		"/",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.CreateAPIKeyPermission),
		apikeys.NewCreateOrganizationAPIKeyHandler(useCases.APIKeys),
	)

	apiKeysRouterGroup.Post(
		"/:apiKeyId/regenerate",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.UpdateAPIKeyPermission),
		apikeys.NewRegenerateOrganizationAPIKeyHandler(useCases.APIKeys),
	)

	apiKeysRouterGroup.Delete(
		"/:apiKeyId",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.DeleteAPIKeyPermission),
		apikeys.NewDeleteOrganizationAPIKeyHandler(useCases.APIKeys),
	)
}

func RegisterAPIKeysRoutes(
	router fiber.Router,
	cacheMiddleware *cachemiddleware.Middleware,
	useCases *apiKeysUseCases.Services,
) {
	apiKeysRouterGroup := router.Group("/api-keys")

	apiKeysRouterGroup.Get(
		"/scopes",
		cacheMiddleware.Cache(1*time.Hour),
		apikeys.NewGetAPIKeyScopesHandler(useCases.APIKeys),
	)
}

func RegisterRoutes(
	router fiber.Router,
	authMiddleware *commonAuth.AuthMiddleware,
	cacheMiddleware *cachemiddleware.Middleware,
	useCases *apiKeysUseCases.Services,
) {
	RegisterOrganizationAPIKeysRoutes(router, authMiddleware, useCases)
	RegisterAPIKeysRoutes(router, cacheMiddleware, useCases)
}
