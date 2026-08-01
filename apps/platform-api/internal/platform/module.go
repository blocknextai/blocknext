package platform

import (
	cachemiddleware "github.com/blocknextai/go-packages/fiber/middleware/cache"
	"github.com/blocknextai/platform-api/internal/common/presentation/auth"
	nodeEngineApplicationCredentials "github.com/blocknextai/platform-api/internal/nodeengine/application/credentials"
	platformApplicationPlatformCredentials "github.com/blocknextai/platform-api/internal/platform/application/platformcredentials"
	platformInfrastructure "github.com/blocknextai/platform-api/internal/platform/infrastructure"
	platformInfrastructurePlatformConfig "github.com/blocknextai/platform-api/internal/platform/infrastructure/platformconfig"
	platformPresentation "github.com/blocknextai/platform-api/internal/platform/presentation"
	"github.com/gofiber/fiber/v3"
)

type Dependencies struct {
	CredentialConfigs map[string]string

	NodeEngineCredentialService nodeEngineApplicationCredentials.CredentialService
}

type Module struct {
	PlatformCredentialService platformApplicationPlatformCredentials.PlatformCredentialService

	handlers *platformInfrastructure.Handlers
}

func NewModule(deps Dependencies) *Module {
	loader := platformInfrastructurePlatformConfig.NewPlatformConfigLoader(
		deps.CredentialConfigs,
		deps.NodeEngineCredentialService.GetAllCredentials(),
	)
	service := platformApplicationPlatformCredentials.NewPlatformCredentialService(loader)
	handlers := platformInfrastructure.RegisterInfrastructure(platformInfrastructure.RegisterInfrastructureDeps{
		PlatformCredentialService: service,
		CredentialService:         deps.NodeEngineCredentialService,
	})
	return &Module{
		PlatformCredentialService: service,
		handlers:                  handlers,
	}
}

func (m *Module) Register(router fiber.Router, authMiddleware *auth.AuthMiddleware, cacheMiddleware *cachemiddleware.Middleware) {
	platformPresentation.RegisterPresentation(router, authMiddleware, cacheMiddleware, m.handlers)
}
