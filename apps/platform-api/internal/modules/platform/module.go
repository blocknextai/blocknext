package platform

import (
	"github.com/gofiber/fiber/v3"

	cachemiddleware "github.com/blocknextai/go-packages/fiber/middleware/cache"
	"github.com/blocknextai/platform-api/internal/common/auth"
	nodeengineContract "github.com/blocknextai/platform-api/internal/modules/nodeengine/contract"
	platformApplicationPlatformCredentials "github.com/blocknextai/platform-api/internal/modules/platform/internal/application/platformcredentials"
	platformHTTP "github.com/blocknextai/platform-api/internal/modules/platform/internal/http"
	platformPlatformConfig "github.com/blocknextai/platform-api/internal/modules/platform/internal/platformconfig"
	platformUseCases "github.com/blocknextai/platform-api/internal/modules/platform/internal/usecases"
)

type Dependencies struct {
	CredentialConfigs map[string]string

	FunctionCallingEnabled     bool
	WorkflowsGenerationEnabled bool

	NodeEngineCredentialService nodeengineContract.CredentialService
}

type Module struct {
	PlatformCredentialService platformApplicationPlatformCredentials.PlatformCredentialService

	useCases *platformUseCases.Services
}

func NewModule(deps Dependencies) *Module {
	loader := platformPlatformConfig.NewPlatformConfigLoader(
		deps.CredentialConfigs,
		deps.NodeEngineCredentialService.GetAllCredentials(),
	)
	service := platformApplicationPlatformCredentials.NewPlatformCredentialService(loader)
	useCases := platformUseCases.NewServices(platformUseCases.ServiceDependencies{
		FunctionCallingEnabled:     deps.FunctionCallingEnabled,
		WorkflowsGenerationEnabled: deps.WorkflowsGenerationEnabled,

		PlatformCredentialService: service,
		CredentialService:         deps.NodeEngineCredentialService,
	})
	return &Module{
		PlatformCredentialService: service,
		useCases:                  useCases,
	}
}

func (m *Module) Register(router fiber.Router, authMiddleware *auth.AuthMiddleware, cacheMiddleware *cachemiddleware.Middleware) {
	platformHTTP.RegisterRoutes(router, authMiddleware, cacheMiddleware, m.useCases)
}
