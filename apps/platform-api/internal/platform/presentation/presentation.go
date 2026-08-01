package presentation

import (
	"time"

	cachemiddleware "github.com/blocknextai/go-packages/fiber/middleware/cache"
	"github.com/blocknextai/go-packages/rbac"
	commonPresentationAuth "github.com/blocknextai/platform-api/internal/common/presentation/auth"
	platformInfrastructure "github.com/blocknextai/platform-api/internal/platform/infrastructure"
	"github.com/blocknextai/platform-api/internal/platform/presentation/platformcredentials"
	"github.com/gofiber/fiber/v3"
)

func RegisterPresentation(
	router fiber.Router,
	authMiddleware *commonPresentationAuth.AuthMiddleware,
	cacheMiddleware *cachemiddleware.Middleware,
	handlers *platformInfrastructure.Handlers,
) {
	platformCredentialsRouterGroup := router.Group("/platform/credentials")

	platformCredentialsRouterGroup.Get(
		"/",
		authMiddleware.Authenticate(),
		authMiddleware.RequireUserPermission(rbac.ReadPlatformCredentialPermission),
		cacheMiddleware.Cache(5*time.Minute),
		platformcredentials.NewGetAllPlatformCredentialsHandler(handlers.GetAllPlatformCredentials),
	)

	platformCredentialsRouterGroup.Get(
		"/:id",
		authMiddleware.Authenticate(),
		authMiddleware.RequireUserPermission(rbac.ReadPlatformCredentialPermission),
		cacheMiddleware.Cache(5*time.Minute),
		platformcredentials.NewGetPlatformCredentialByIDHandler(handlers.GetPlatformCredentialByID),
	)
}
