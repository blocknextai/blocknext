package apikeys

import (
	"database/sql"

	commonAuth "github.com/blocknextai/platform-api/internal/common/auth"

	"github.com/gofiber/fiber/v3"

	"github.com/blocknextai/go-packages/database"
	cachemiddleware "github.com/blocknextai/go-packages/fiber/middleware/cache"
	apiKeysApplicationAPIKeys "github.com/blocknextai/platform-api/internal/modules/apikeys/internal/application/apikeys"
	apiKeysDomainAPIKeys "github.com/blocknextai/platform-api/internal/modules/apikeys/internal/domain/apikeys"
	apiKeysHTTP "github.com/blocknextai/platform-api/internal/modules/apikeys/internal/http"
	apiKeysPostgresAPIKeys "github.com/blocknextai/platform-api/internal/modules/apikeys/internal/postgres/apikeys"
	apiKeysUseCases "github.com/blocknextai/platform-api/internal/modules/apikeys/internal/usecases"
)

type Dependencies struct {
	DB                 *sql.DB
	TransactionManager database.TransactionManager
}

type Module struct {
	APIKeyValidator commonAuth.APIKeyValidator[apiKeysDomainAPIKeys.Scope]
	APIKeyService   apiKeysApplicationAPIKeys.APIKeyService

	useCases *apiKeysUseCases.Services
}

func NewModule(deps Dependencies) *Module {
	repository := apiKeysPostgresAPIKeys.NewAPIKeyRepository(deps.DB)
	apiKeyValidator := apiKeysApplicationAPIKeys.NewAPIKeyValidator(repository)
	apiKeyService := apiKeysApplicationAPIKeys.NewAPIKeyService(repository)
	useCases := apiKeysUseCases.NewServices(apiKeysUseCases.ServiceDependencies{
		TransactionManager: deps.TransactionManager,

		ApiKeyRepository: repository,
	})
	return &Module{
		APIKeyValidator: apiKeyValidator,
		APIKeyService:   apiKeyService,
		useCases:        useCases,
	}
}

func (m *Module) Register(router fiber.Router, authMiddleware *commonAuth.AuthMiddleware, cacheMiddleware *cachemiddleware.Middleware) {
	apiKeysHTTP.RegisterRoutes(router, authMiddleware, cacheMiddleware, m.useCases)
}
