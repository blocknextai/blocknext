package organizations

import (
	"database/sql"

	"github.com/blocknextai/go-packages/cache"
	"github.com/blocknextai/go-packages/database"
	cachemiddleware "github.com/blocknextai/go-packages/fiber/middleware/cache"
	accountApplicationLinkedAccounts "github.com/blocknextai/platform-api/internal/account/application/linkedaccounts"
	accountApplicationUsers "github.com/blocknextai/platform-api/internal/account/application/users"
	commonApplicationAuth "github.com/blocknextai/platform-api/internal/common/application/auth"
	"github.com/blocknextai/platform-api/internal/common/presentation/auth"
	"github.com/blocknextai/platform-api/internal/eventbus/application/publishing"
	organizationsApplicationAuth "github.com/blocknextai/platform-api/internal/organizations/application/auth"
	organizationsApplicationOrganizations "github.com/blocknextai/platform-api/internal/organizations/application/organizations"
	organizationsApplicationOrganizationUsers "github.com/blocknextai/platform-api/internal/organizations/application/organizationusers"
	organizationsInfrastructure "github.com/blocknextai/platform-api/internal/organizations/infrastructure"
	organizationsInfrastructureOrganizations "github.com/blocknextai/platform-api/internal/organizations/infrastructure/organizations"
	organizationsInfrastructureOrganizationUsers "github.com/blocknextai/platform-api/internal/organizations/infrastructure/organizationusers"
	organizationsPresentation "github.com/blocknextai/platform-api/internal/organizations/presentation"
	"github.com/gofiber/fiber/v3"
)

type Dependencies struct {
	DB                       *sql.DB
	TransactionManager       database.TransactionManager
	EventBusPublisherService publishing.PublisherService
	CacheService             cache.Service

	UserService          accountApplicationUsers.UserService
	LinkedAccountService accountApplicationLinkedAccounts.LinkedAccountService
}

type Module struct {
	OrganizationService           organizationsApplicationOrganizations.OrganizationService
	OrganizationUserService       organizationsApplicationOrganizationUsers.OrganizationUserService
	OrganizationPermissionChecker commonApplicationAuth.OrganizationPermissionChecker

	handlers *organizationsInfrastructure.Handlers
}

func NewModule(deps Dependencies) *Module {
	organizationRepository := organizationsInfrastructureOrganizations.NewOrganizationRepository(deps.DB)
	organizationUserRepository := organizationsInfrastructureOrganizationUsers.NewOrganizationUserRepository(deps.DB)

	organizationService := organizationsApplicationOrganizations.NewOrganizationService(organizationRepository)
	organizationUserService := organizationsApplicationOrganizationUsers.NewOrganizationUserService(organizationUserRepository)
	organizationPermissionChecker := organizationsApplicationAuth.NewOrganizationPermissionChecker(organizationUserRepository, deps.CacheService)

	handlers := organizationsInfrastructure.RegisterInfrastructure(organizationsInfrastructure.RegisterInfrastructureDeps{
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
		handlers:                      handlers,
	}
}

func (m *Module) Register(router fiber.Router, authMiddleware *auth.AuthMiddleware, cacheMiddleware *cachemiddleware.Middleware) {
	organizationsPresentation.RegisterPresentation(router, authMiddleware, cacheMiddleware, m.handlers)
}
