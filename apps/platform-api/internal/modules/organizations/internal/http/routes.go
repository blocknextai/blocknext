package http

import (
	"time"

	"github.com/gofiber/fiber/v3"

	cachemiddleware "github.com/blocknextai/go-packages/fiber/middleware/cache"
	"github.com/blocknextai/go-packages/rbac"
	commonAuth "github.com/blocknextai/platform-api/internal/common/auth"
	organizationsHTTPOrganizations "github.com/blocknextai/platform-api/internal/modules/organizations/internal/http/organizations"
	organizationsHTTPOrganizationUsers "github.com/blocknextai/platform-api/internal/modules/organizations/internal/http/organizationusers"
	organizationsUseCases "github.com/blocknextai/platform-api/internal/modules/organizations/internal/usecases"
)

func RegisterRoutes(
	router fiber.Router,
	authMiddleware *commonAuth.AuthMiddleware,
	cacheMiddleware *cachemiddleware.Middleware,
	useCases *organizationsUseCases.Services,
) {
	organizationsRouterGroup := router.Group("/organizations")

	organizationsRouterGroup.Get(
		"/",
		authMiddleware.Authenticate(),
		authMiddleware.RequireUserPermission(rbac.ReadOrganizationPermission),
		organizationsHTTPOrganizations.NewGetAllOrganizationsHandler(useCases.Organizations),
	)

	organizationsRouterGroup.Get(
		"/roles",
		cacheMiddleware.Cache(5*time.Minute),
		organizationsHTTPOrganizationUsers.NewGetOrganizationRolesHandler(useCases.OrganizationUsers),
	)

	organizationsRouterGroup.Post(
		"/",
		authMiddleware.Authenticate(),
		authMiddleware.RequireUserPermission(rbac.CreateOrganizationPermission),
		organizationsHTTPOrganizations.NewCreateOrganizationHandler(useCases.Organizations),
	)

	organizationsRouterGroup.Get(
		"/:organizationId",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.ReadOrganizationPermission),
		organizationsHTTPOrganizations.NewGetOrganizationByIDHandler(useCases.Organizations),
	)

	organizationsRouterGroup.Get(
		"/:organizationId/me",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.ReadOrganizationPermission),
		organizationsHTTPOrganizationUsers.NewGetOrganizationMeHandler(useCases.OrganizationUsers),
	)

	organizationsRouterGroup.Put(
		"/:organizationId",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.UpdateOrganizationPermission),
		organizationsHTTPOrganizations.NewUpdateOrganizationHandler(useCases.Organizations),
	)

	organizationsRouterGroup.Delete(
		"/:organizationId",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.DeleteOrganizationPermission),
		organizationsHTTPOrganizations.NewDeleteOrganizationHandler(useCases.Organizations),
	)

	organizationUsersRouterGroup := organizationsRouterGroup.Group("/:organizationId/users")

	organizationUsersRouterGroup.Get(
		"/",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.ReadOrganizationUserPermission),
		organizationsHTTPOrganizationUsers.NewGetAllOrganizationUsersHandler(useCases.OrganizationUsers),
	)

	organizationUsersRouterGroup.Post(
		"/",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.CreateOrganizationUserPermission),
		organizationsHTTPOrganizationUsers.NewCreateOrganizationUserHandler(useCases.OrganizationUsers),
	)

	organizationUsersRouterGroup.Get(
		"/:userId",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.ReadOrganizationUserPermission),
		organizationsHTTPOrganizationUsers.NewGetOrganizationUserByUserIDHandler(useCases.OrganizationUsers),
	)

	organizationUsersRouterGroup.Put(
		"/:userId/info",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.UpdateOrganizationUserInfoPermission),
		organizationsHTTPOrganizationUsers.NewUpdateOrganizationUserInfoHandler(useCases.OrganizationUsers),
	)

	organizationUsersRouterGroup.Put(
		"/:userId/role",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.UpdateOrganizationUserRolePermission),
		organizationsHTTPOrganizationUsers.NewUpdateOrganizationUserRoleHandler(useCases.OrganizationUsers),
	)

	organizationUsersRouterGroup.Delete(
		"/:userId",
		authMiddleware.Authenticate(),
		authMiddleware.RequireOrganizationPermission(rbac.DeleteOrganizationUserPermission),
		organizationsHTTPOrganizationUsers.NewDeleteOrganizationUserHandler(useCases.OrganizationUsers),
	)

}
