package presentation

import (
	"github.com/blocknextai/go-packages/rbac"
	commonPresentationAuth "github.com/blocknextai/platform-api/internal/common/presentation/auth"
	credentialsInfrastructure "github.com/blocknextai/platform-api/internal/credentials/infrastructure"
	"github.com/blocknextai/platform-api/internal/credentials/presentation/credentials"
	"github.com/gofiber/fiber/v3"
)

func RegisterOrganizationCredentialsPresentation(
	router fiber.Router,
	authMiddleware *commonPresentationAuth.AuthMiddleware,
	handlers *credentialsInfrastructure.Handlers,
) {
	credentialsRouterGroup := router.Group("/organizations/:organizationId/credentials")

	credentialsRouterGroup.Get(
		"/",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.ReadOrganizationCredentialsPermission),
		credentials.NewGetAllOrganizationCredentialsHandler(handlers.GetAllCredentials),
	)

	credentialsRouterGroup.Get(
		"/by-nodes",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.ReadOrganizationCredentialsPermission),
		credentials.NewGetOrganizationCredentialsForNodesHandler(handlers.GetCredentialsForNodes),
	)

	credentialsRouterGroup.Get(
		"/:credentialId",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.ReadOrganizationCredentialsPermission),
		credentials.NewGetOrganizationCredentialByIDHandler(handlers.GetCredentialByID),
	)

	credentialsRouterGroup.Post(
		"/",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.CreateOrganizationCredentialsPermission),
		credentials.NewCreateOrganizationCredentialHandler(handlers.CreateCredential),
	)

	credentialsRouterGroup.Put(
		"/:credentialId",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.UpdateOrganizationCredentialsPermission),
		credentials.NewUpdateOrganizationCredentialHandler(handlers.UpdateCredential),
	)

	credentialsRouterGroup.Delete(
		"/:credentialId",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.DeleteOrganizationCredentialsPermission),
		credentials.NewDeleteOrganizationCredentialHandler(handlers.DeleteCredential),
	)
}

func RegisterPresentation(
	router fiber.Router,
	authMiddleware *commonPresentationAuth.AuthMiddleware,
	handlers *credentialsInfrastructure.Handlers,
) {
	RegisterOrganizationCredentialsPresentation(router, authMiddleware, handlers)
}
