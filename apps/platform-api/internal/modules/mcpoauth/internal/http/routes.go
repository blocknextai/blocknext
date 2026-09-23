package http

import (
	"time"

	"github.com/gofiber/fiber/v3"

	cachemiddleware "github.com/blocknextai/go-packages/fiber/middleware/cache"
	"github.com/blocknextai/go-packages/rbac"
	commonAuth "github.com/blocknextai/platform-api/internal/common/auth"
	mcpOAuthDomainOAuth2 "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/domain/oauth2"
	"github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/http/authorizationrequests"
	"github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/http/clients"
	"github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/http/grants"
	"github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/http/metadata"
	"github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/http/tokens"
	mcpOAuthUseCases "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/usecases"
)

func RegisterMetadataRoutes(
	router fiber.Router,
	cacheMiddleware *cachemiddleware.Middleware,
	useCases *mcpOAuthUseCases.Services,
) {
	router.Get(
		mcpOAuthDomainOAuth2.AuthorizationServerMetadataPath,
		cacheMiddleware.Cache(1*time.Hour),
		metadata.NewGetAuthorizationServerMetadataHandler(useCases.Metadata),
	)
}

func RegisterClientsRoutes(
	router fiber.Router,
	useCases *mcpOAuthUseCases.Services,
) {
	router.Post(
		mcpOAuthDomainOAuth2.RegistrationEndpointPath,
		clients.NewRegisterClientHandler(useCases.Clients),
	)
}

func RegisterAuthorizationRequestsRoutes(
	router fiber.Router,
	authMiddleware *commonAuth.AuthMiddleware,
	useCases *mcpOAuthUseCases.Services,
) {
	router.Get(
		mcpOAuthDomainOAuth2.AuthorizationEndpointPath,
		authorizationrequests.NewCreateAuthorizationRequestHandler(useCases.AuthorizationRequests),
	)

	authorizationRequestsRouterGroup := router.Group("/mcp-oauth/authorization-requests")

	authorizationRequestsRouterGroup.Get(
		"/:authorizationRequestId",
		authMiddleware.Authenticate(),
		authorizationrequests.NewGetAuthorizationRequestHandler(useCases.AuthorizationRequests),
	)

	authorizationRequestsRouterGroup.Post(
		"/:authorizationRequestId/approve",
		authMiddleware.Authenticate(),
		authorizationrequests.NewApproveAuthorizationRequestHandler(useCases.AuthorizationRequests),
	)

	authorizationRequestsRouterGroup.Post(
		"/:authorizationRequestId/deny",
		authMiddleware.Authenticate(),
		authorizationrequests.NewDenyAuthorizationRequestHandler(useCases.AuthorizationRequests),
	)
}

func RegisterTokensRoutes(
	router fiber.Router,
	useCases *mcpOAuthUseCases.Services,
) {
	router.Post(
		mcpOAuthDomainOAuth2.TokenEndpointPath,
		tokens.NewExchangeTokenHandler(useCases.Tokens),
	)

	router.Post(
		mcpOAuthDomainOAuth2.RevocationEndpointPath,
		tokens.NewRevokeTokenHandler(useCases.Tokens),
	)
}

func RegisterUserGrantsRoutes(
	router fiber.Router,
	authMiddleware *commonAuth.AuthMiddleware,
	useCases *mcpOAuthUseCases.Services,
) {
	grantsRouterGroup := router.Group("/users/me/mcp-oauth/grants")

	grantsRouterGroup.Get(
		"/",
		authMiddleware.Authenticate(),
		authMiddleware.RequireUserPermission(rbac.ReadUserPermission),
		grants.NewGetAllUserGrantsHandler(useCases.Grants),
	)

	grantsRouterGroup.Delete(
		"/:grantId",
		authMiddleware.Authenticate(),
		authMiddleware.RequireUserPermission(rbac.UpdateUserPermission),
		grants.NewRevokeUserGrantHandler(useCases.Grants),
	)
}

func RegisterRoutes(
	router fiber.Router,
	authMiddleware *commonAuth.AuthMiddleware,
	cacheMiddleware *cachemiddleware.Middleware,
	useCases *mcpOAuthUseCases.Services,
) {
	RegisterMetadataRoutes(router, cacheMiddleware, useCases)
	RegisterClientsRoutes(router, useCases)
	RegisterAuthorizationRequestsRoutes(router, authMiddleware, useCases)
	RegisterTokensRoutes(router, useCases)
	RegisterUserGrantsRoutes(router, authMiddleware, useCases)
}
