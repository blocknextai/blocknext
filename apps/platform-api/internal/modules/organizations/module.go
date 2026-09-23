package organizations

import (
	"database/sql"

	commonAuth "github.com/blocknextai/platform-api/internal/common/auth"

	"github.com/gofiber/fiber/v3"

	"github.com/blocknextai/go-packages/cache"
	"github.com/blocknextai/go-packages/database"
	cachemiddleware "github.com/blocknextai/go-packages/fiber/middleware/cache"
	"github.com/blocknextai/platform-api/internal/eventbus/publishing"
	accountContract "github.com/blocknextai/platform-api/internal/modules/account/contract"
	organizationsApplicationAuth "github.com/blocknextai/platform-api/internal/modules/organizations/internal/application/auth"
	organizationsApplicationOrganizations "github.com/blocknextai/platform-api/internal/modules/organizations/internal/application/organizations"
	organizationsApplicationOrganizationUsers "github.com/blocknextai/platform-api/internal/modules/organizations/internal/application/organizationusers"
	organizationsHTTP "github.com/blocknextai/platform-api/internal/modules/organizations/internal/http"
	organizationsPostgresOrganizations "github.com/blocknextai/platform-api/internal/modules/organizations/internal/postgres/organizations"
	organizationsPostgresOrganizationUsers "github.com/blocknextai/platform-api/internal/modules/organizations/internal/postgres/organizationusers"
	organizationsUseCases "github.com/blocknextai/platform-api/internal/modules/organizations/internal/usecases"
)

type Dependencies struct {
	DB                       *sql.DB
	TransactionManager       database.TransactionManager
	EventBusPublisherService publishing.PublisherService
	CacheService             cache.Service

	UserService          accountContract.UserService
	LinkedAccountService accountContract.LinkedAccountService
}

type Module struct {
	OrganizationService           organizationsApplicationOrganizations.OrganizationService
	OrganizationUserService       organizationsApplicationOrganizationUsers.OrganizationUserService
	OrganizationPermissionChecker commonAuth.OrganizationPermissionChecker

	useCases *organizationsUseCases.Services
}

func NewModule(deps Dependencies) *Module {
	organizationRepository := organizationsPostgresOrganizations.NewOrganizationRepository(deps.DB)
	organizationUserRepository := organizationsPostgresOrganizationUsers.NewOrganizationUserRepository(deps.DB)

	organizationService := organizationsApplicationOrganizations.NewOrganizationService(organizationRepository)
	organizationUserService := organizationsApplicationOrganizationUsers.NewOrganizationUserService(organizationUserRepository)
	organizationPermissionChecker := organizationsApplicationAuth.NewOrganizationPermissionChecker(organizationUserRepository, deps.CacheService)

	useCases := organizationsUseCases.NewServices(organizationsUseCases.ServiceDependencies{
		TransactionManager:       deps.TransactionManager,
		EventBusPublisherService: deps.EventBusPublisherService,

		OrganizationRepository:     organizationRepository,
		OrganizationUserRepository: organizationUserRepository,
		UserService:                deps.UserService,
		LinkedAccountService:       deps.LinkedAccountService,
	})

	return &Module{
		OrganizationService:           organizationService,
		OrganizationUserService:       organizationUserService,
		OrganizationPermissionChecker: organizationPermissionChecker,
		useCases:                      useCases,
	}
}

func (m *Module) Register(router fiber.Router, authMiddleware *commonAuth.AuthMiddleware, cacheMiddleware *cachemiddleware.Middleware) {
	organizationsHTTP.RegisterRoutes(router, authMiddleware, cacheMiddleware, m.useCases)
}
