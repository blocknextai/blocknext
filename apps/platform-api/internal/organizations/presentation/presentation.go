package presentation

import (
	"time"

	"github.com/gofiber/fiber/v3"

	cachemiddleware "github.com/blocknextai/go-packages/fiber/middleware/cache"
	"github.com/blocknextai/go-packages/rbac"
	commonPresentationAuth "github.com/blocknextai/platform-api/internal/common/presentation/auth"
	organizationsInfrastructure "github.com/blocknextai/platform-api/internal/organizations/infrastructure"
	organizationsPresentationOrganizations "github.com/blocknextai/platform-api/internal/organizations/presentation/organizations"
	organizationsPresentationOrganizationUsers "github.com/blocknextai/platform-api/internal/organizations/presentation/organizationusers"
)

func RegisterPresentation(
	router fiber.Router,
	authMiddleware *commonPresentationAuth.AuthMiddleware,
	cacheMiddleware *cachemiddleware.Middleware,
	handlers *organizationsInfrastructure.Handlers,
) {
	organizationsRouterGroup := router.Group("/organizations")

	organizationsRouterGroup.Get(
		"/",
		authMiddleware.Authenticate(),
		authMiddleware.RequireUserPermission(rbac.ReadOrganizationPermission),
		organizationsPresentationOrganizations.NewGetAllOrganizationsHandler(handlers.GetAllOrganizations),
	)

	organizationsRouterGroup.Get(
		"/roles",
		cacheMiddleware.Cache(5*time.Minute),
		organizationsPresentationOrganizationUsers.NewGetOrganizationRolesHandler(handlers.GetRoles),
	)

	organizationsRouterGroup.Post(
		"/",
		authMiddleware.Authenticate(),
		authMiddleware.RequireUserPermission(rbac.CreateOrganizationPermission),
		organizationsPresentationOrganizations.NewCreateOrganizationHandler(handlers.CreateOrganization),
	)

	organizationsRouterGroup.Get(
		"/:organizationId",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.ReadOrganizationPermission),
		organizationsPresentationOrganizations.NewGetOrganizationByIDHandler(handlers.GetOrganizationByID),
	)

	organizationsRouterGroup.Get(
		"/:organizationId/me",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.ReadOrganizationPermission),
		organizationsPresentationOrganizationUsers.NewGetOrganizationMeHandler(handlers.GetOrganizationMe),
	)

	organizationsRouterGroup.Put(
		"/:organizationId",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.UpdateOrganizationPermission),
		organizationsPresentationOrganizations.NewUpdateOrganizationHandler(handlers.UpdateOrganization),
	)

	organizationsRouterGroup.Delete(
		"/:organizationId",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.DeleteOrganizationPermission),
		organizationsPresentationOrganizations.NewDeleteOrganizationHandler(handlers.DeleteOrganization),
	)

	organizationUsersRouterGroup := organizationsRouterGroup.Group("/:organizationId/users")

	organizationUsersRouterGroup.Get(
		"/",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.ReadOrganizationUserPermission),
		organizationsPresentationOrganizationUsers.NewGetAllOrganizationUsersHandler(handlers.GetAllOrganizationUsers),
	)

	organizationUsersRouterGroup.Post(
		"/",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.CreateOrganizationUserPermission),
		organizationsPresentationOrganizationUsers.NewCreateOrganizationUserHandler(handlers.CreateOrganizationUser),
	)

	organizationUsersRouterGroup.Get(
		"/:userId",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.ReadOrganizationUserPermission),
		organizationsPresentationOrganizationUsers.NewGetOrganizationUserByUserIDHandler(handlers.GetOrganizationUserByUserID),
	)

	organizationUsersRouterGroup.Put(
		"/:userId/info",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.UpdateOrganizationUserInfoPermission),
		organizationsPresentationOrganizationUsers.NewUpdateOrganizationUserInfoHandler(handlers.UpdateOrganizationUserInfo),
	)

	organizationUsersRouterGroup.Put(
		"/:userId/role",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.UpdateOrganizationUserRolePermission),
		organizationsPresentationOrganizationUsers.NewUpdateOrganizationUserRoleHandler(handlers.UpdateOrganizationUserRole),
	)

	organizationUsersRouterGroup.Delete(
		"/:userId",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.DeleteOrganizationUserPermission),
		organizationsPresentationOrganizationUsers.NewDeleteOrganizationUserHandler(handlers.DeleteOrganizationUser),
	)

}
