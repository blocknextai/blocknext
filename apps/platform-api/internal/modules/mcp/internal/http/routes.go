package http

import (
	"slices"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/adaptor"

	cachemiddleware "github.com/blocknextai/go-packages/fiber/middleware/cache"
	commonAuth "github.com/blocknextai/platform-api/internal/common/auth"
	apikeysContract "github.com/blocknextai/platform-api/internal/modules/apikeys/contract"
	"github.com/blocknextai/platform-api/internal/modules/mcp/internal/application/servers"
	mcpUseCases "github.com/blocknextai/platform-api/internal/modules/mcp/internal/usecases"
)

const (
	protectedResourceMetadataPath = "/.well-known/oauth-protected-resource"
)

func RegisterRoutes(
	router fiber.Router,
	cacheMiddleware *cachemiddleware.Middleware,
	useCases *mcpUseCases.Services,
	apiKeyMiddleware *commonAuth.APIKeyMiddleware[apikeysContract.Scope],
	accessTokenMiddleware *commonAuth.AccessTokenMiddleware,
) {
	router.Get(
		"/servers",
		cacheMiddleware.Cache(5*time.Minute),
		NewGetAllServersHandler(useCases.Servers),
	)

	router.Get(
		protectedResourceMetadataPath+"/:serverId/mcp",
		NewGetProtectedResourceMetadataHandler(useCases.Servers),
	)

	for _, serverHandler := range useCases.ServerHandlers {
		handler := adaptor.HTTPHandler(serverHandler.Handler)
		path := "/" + serverHandler.ID + "/mcp"
		requireAuthentication := authenticate(serverHandler.AuthMethods, apiKeyMiddleware, accessTokenMiddleware)
		requireServerScope := requireScope(serverHandler.AuthMethods, serverHandler.Scopes, apiKeyMiddleware, accessTokenMiddleware)

		router.Post(path, requireAuthentication, requireServerScope, handler)
		router.Get(path, requireAuthentication, requireServerScope, handler)
		router.Delete(path, requireAuthentication, requireServerScope, handler)
	}
}

func authenticate(
	authMethods []servers.AuthMethod,
	apiKeyMiddleware *commonAuth.APIKeyMiddleware[apikeysContract.Scope],
	accessTokenMiddleware *commonAuth.AccessTokenMiddleware,
) fiber.Handler {
	accessTokenAuthentication := accessTokenMiddleware.Authenticate()
	if !slices.Contains(authMethods, servers.AuthMethodAPIKey) {
		return accessTokenAuthentication
	}

	apiKeyAuthentication := apiKeyMiddleware.Authenticate()

	return func(c fiber.Ctx) error {
		if commonAuth.IsAccessTokenRequest(c) || c.Get(commonAuth.APIKeyHeader) == "" {
			return accessTokenAuthentication(c)
		}

		return apiKeyAuthentication(c)
	}
}

func requireScope(
	authMethods []servers.AuthMethod,
	scopes []string,
	apiKeyMiddleware *commonAuth.APIKeyMiddleware[apikeysContract.Scope],
	accessTokenMiddleware *commonAuth.AccessTokenMiddleware,
) fiber.Handler {
	accessTokenScope := accessTokenMiddleware.RequireAnyScope(scopes...)
	if !slices.Contains(authMethods, servers.AuthMethodAPIKey) {
		return accessTokenScope
	}

	apiKeyScope := apiKeyMiddleware.RequireScope(apikeysContract.ScopeMCPInvoke)

	return func(c fiber.Ctx) error {
		if commonAuth.IsAccessTokenRequest(c) {
			return accessTokenScope(c)
		}

		return apiKeyScope(c)
	}
}
