package http

import (
	"github.com/gofiber/fiber/v3"

	"github.com/blocknextai/go-packages/rbac"
	commonAuth "github.com/blocknextai/platform-api/internal/common/auth"
	"github.com/blocknextai/platform-api/internal/modules/credentials/internal/http/credentials"
	credentialsUseCases "github.com/blocknextai/platform-api/internal/modules/credentials/internal/usecases"
)

func RegisterOrganizationCredentialsRoutes(
	router fiber.Router,
	authMiddleware *commonAuth.AuthMiddleware,
	useCases *credentialsUseCases.Services,
) {
	credentialsRouterGroup := router.Group("/organizations/:organizationId/credentials")

	credentialsRouterGroup.Get(
		"/",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.ReadOrganizationCredentialsPermission),
		credentials.NewGetAllOrganizationCredentialsHandler(useCases.Credentials),
	)

	credentialsRouterGroup.Get(
		"/by-nodes",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.ReadOrganizationCredentialsPermission),
		credentials.NewGetOrganizationCredentialsForNodesHandler(useCases.Credentials),
	)

	credentialsRouterGroup.Get(
		"/:credentialId",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.ReadOrganizationCredentialsPermission),
		credentials.NewGetOrganizationCredentialByIDHandler(useCases.Credentials),
	)

	credentialsRouterGroup.Post(
		"/",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.CreateOrganizationCredentialsPermission),
		credentials.NewCreateOrganizationCredentialHandler(useCases.Credentials),
	)

	credentialsRouterGroup.Put(
		"/:credentialId",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.UpdateOrganizationCredentialsPermission),
		credentials.NewUpdateOrganizationCredentialHandler(useCases.Credentials),
	)

	credentialsRouterGroup.Delete(
		"/:credentialId",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.DeleteOrganizationCredentialsPermission),
		credentials.NewDeleteOrganizationCredentialHandler(useCases.Credentials),
	)
}

func RegisterRoutes(
	router fiber.Router,
	authMiddleware *commonAuth.AuthMiddleware,
	useCases *credentialsUseCases.Services,
) {
	RegisterOrganizationCredentialsRoutes(router, authMiddleware, useCases)
}
