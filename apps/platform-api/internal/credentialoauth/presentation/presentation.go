package presentation

import (
	"github.com/blocknextai/go-packages/rbac"
	commonPresentationAuth "github.com/blocknextai/platform-api/internal/common/presentation/auth"
	credentialOAuthInfrastructure "github.com/blocknextai/platform-api/internal/credentialoauth/infrastructure"
	"github.com/blocknextai/platform-api/internal/credentialoauth/presentation/oauth2"
	"github.com/gofiber/fiber/v3"
)

func RegisterOrganizationCredentialOAuthPresentation(
	router fiber.Router,
	authMiddleware *commonPresentationAuth.AuthMiddleware,
	handlers *credentialOAuthInfrastructure.Handlers,
) {
	credentialOAuthRouterGroup := router.Group("/organizations/:organizationId/credential-oauth/oauth2")

	credentialOAuthRouterGroup.Post(
		"/auth",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.UpdateOrganizationCredentialsPermission),
		oauth2.NewOrganizationAuthHandler(handlers.AuthURL),
	)
}

func RegisterCallbackPresentation(
	router fiber.Router,
	handlers *credentialOAuthInfrastructure.Handlers,
) {
	callbackRouterGroup := router.Group("/credential-oauth/oauth2")

	callbackRouterGroup.Get(
		"/callback",
		oauth2.NewCallbackHandler(handlers.ExchangeCode),
	)
}

func RegisterPresentation(
	router fiber.Router,
	authMiddleware *commonPresentationAuth.AuthMiddleware,
	handlers *credentialOAuthInfrastructure.Handlers,
) {
	RegisterOrganizationCredentialOAuthPresentation(router, authMiddleware, handlers)
	RegisterCallbackPresentation(router, handlers)
}
