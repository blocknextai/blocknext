package http

import (
	"github.com/gofiber/fiber/v3"

	"github.com/blocknextai/go-packages/rbac"
	commonAuth "github.com/blocknextai/platform-api/internal/common/auth"
	"github.com/blocknextai/platform-api/internal/modules/credentialoauth/internal/http/oauth2"
	credentialOAuthUseCases "github.com/blocknextai/platform-api/internal/modules/credentialoauth/internal/usecases"
)

func RegisterOrganizationCredentialOAuthRoutes(
	router fiber.Router,
	authMiddleware *commonAuth.AuthMiddleware,
	useCases *credentialOAuthUseCases.Services,
) {
	credentialOAuthRouterGroup := router.Group("/organizations/:organizationId/credential-oauth/oauth2")

	credentialOAuthRouterGroup.Post(
		"/auth",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.UpdateOrganizationCredentialsPermission),
		oauth2.NewOrganizationAuthHandler(useCases.OAuth2),
	)
}

func RegisterCallbackRoutes(
	router fiber.Router,
	useCases *credentialOAuthUseCases.Services,
) {
	callbackRouterGroup := router.Group("/credential-oauth/oauth2")

	callbackRouterGroup.Get(
		"/callback",
		oauth2.NewCallbackHandler(useCases.OAuth2),
	)
}

func RegisterRoutes(
	router fiber.Router,
	authMiddleware *commonAuth.AuthMiddleware,
	useCases *credentialOAuthUseCases.Services,
) {
	RegisterOrganizationCredentialOAuthRoutes(router, authMiddleware, useCases)
	RegisterCallbackRoutes(router, useCases)
}
