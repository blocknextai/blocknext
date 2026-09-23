package http

import (
	"time"

	"github.com/gofiber/fiber/v3"

	cachemiddleware "github.com/blocknextai/go-packages/fiber/middleware/cache"
	"github.com/blocknextai/go-packages/rbac"
	commonAuth "github.com/blocknextai/platform-api/internal/common/auth"
	"github.com/blocknextai/platform-api/internal/modules/platform/internal/http/features"
	"github.com/blocknextai/platform-api/internal/modules/platform/internal/http/platformcredentials"
	platformUseCases "github.com/blocknextai/platform-api/internal/modules/platform/internal/usecases"
)

func RegisterRoutes(
	router fiber.Router,
	authMiddleware *commonAuth.AuthMiddleware,
	cacheMiddleware *cachemiddleware.Middleware,
	useCases *platformUseCases.Services,
) {
	platformRouterGroup := router.Group("/platform")

	platformRouterGroup.Get(
		"/features",
		cacheMiddleware.Cache(5*time.Minute),
		features.NewGetFeaturesHandler(useCases.Features),
	)

	platformCredentialsRouterGroup := platformRouterGroup.Group("/credentials")

	platformCredentialsRouterGroup.Get(
		"/",
		authMiddleware.Authenticate(),
		authMiddleware.RequireUserPermission(rbac.ReadPlatformCredentialPermission),
		cacheMiddleware.Cache(5*time.Minute),
		platformcredentials.NewGetAllPlatformCredentialsHandler(useCases.PlatformCredentials),
	)

	platformCredentialsRouterGroup.Get(
		"/:id",
		authMiddleware.Authenticate(),
		authMiddleware.RequireUserPermission(rbac.ReadPlatformCredentialPermission),
		cacheMiddleware.Cache(5*time.Minute),
		platformcredentials.NewGetPlatformCredentialByIDHandler(useCases.PlatformCredentials),
	)
}
