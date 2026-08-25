package apikeys

import (
	"database/sql"

	"github.com/blocknextai/go-packages/database"
	cachemiddleware "github.com/blocknextai/go-packages/fiber/middleware/cache"
	apiKeysApplicationAPIKeys "github.com/blocknextai/platform-api/internal/apikeys/application/apikeys"
	apiKeysInfrastructure "github.com/blocknextai/platform-api/internal/apikeys/infrastructure"
	apiKeysInfrastructureAPIKeys "github.com/blocknextai/platform-api/internal/apikeys/infrastructure/apikeys"
	apiKeysPresentation "github.com/blocknextai/platform-api/internal/apikeys/presentation"
	"github.com/blocknextai/platform-api/internal/common/presentation/auth"
	"github.com/gofiber/fiber/v3"
)

type Dependencies struct {
	DB                 *sql.DB
	TransactionManager database.TransactionManager
}

type Module struct {
	APIKeyValidator apiKeysApplicationAPIKeys.APIKeyValidator
	APIKeyService   apiKeysApplicationAPIKeys.APIKeyService

	handlers *apiKeysInfrastructure.Handlers
}

func NewModule(deps Dependencies) *Module {
	repository := apiKeysInfrastructureAPIKeys.NewAPIKeyRepository(deps.DB)
	apiKeyValidator := apiKeysApplicationAPIKeys.NewAPIKeyValidator(repository)
	apiKeyService := apiKeysApplicationAPIKeys.NewAPIKeyService(repository)
	handlers := apiKeysInfrastructure.RegisterInfrastructure(apiKeysInfrastructure.RegisterInfrastructureDeps{
		TransactionManager: deps.TransactionManager,

		ApiKeyRepository: repository,
	})
	return &Module{
		APIKeyValidator: apiKeyValidator,
		APIKeyService:   apiKeyService,
		handlers:        handlers,
	}
}

func (m *Module) Register(router fiber.Router, authMiddleware *auth.AuthMiddleware, cacheMiddleware *cachemiddleware.Middleware) {
	apiKeysPresentation.RegisterPresentation(router, authMiddleware, cacheMiddleware, m.handlers)
}
