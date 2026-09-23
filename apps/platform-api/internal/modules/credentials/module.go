package credentials

import (
	"database/sql"

	"github.com/gofiber/fiber/v3"

	"github.com/blocknextai/go-packages/database"
	"github.com/blocknextai/go-packages/secretmanager"
	"github.com/blocknextai/platform-api/internal/common/auth"
	credentialsApplicationCredentials "github.com/blocknextai/platform-api/internal/modules/credentials/internal/application/credentials"
	credentialsHTTP "github.com/blocknextai/platform-api/internal/modules/credentials/internal/http"
	credentialsPostgresCredentials "github.com/blocknextai/platform-api/internal/modules/credentials/internal/postgres/credentials"
	credentialsUseCases "github.com/blocknextai/platform-api/internal/modules/credentials/internal/usecases"
	nodeengineContract "github.com/blocknextai/platform-api/internal/modules/nodeengine/contract"
	platformContract "github.com/blocknextai/platform-api/internal/modules/platform/contract"
)

type Dependencies struct {
	DB                 *sql.DB
	TransactionManager database.TransactionManager
	SecretManager      secretmanager.SecretManager

	NodeEngineCredentialProcessor nodeengineContract.CredentialProcessor
	NodeEngineCredentialService   nodeengineContract.CredentialService
	PlatformCredentialService     platformContract.PlatformCredentialService
}

type Module struct {
	CredentialService credentialsApplicationCredentials.CredentialService

	useCases *credentialsUseCases.Services
}

func NewModule(deps Dependencies) *Module {
	repository := credentialsPostgresCredentials.NewCredentialRepository(deps.DB)
	service := credentialsApplicationCredentials.NewCredentialService(repository, deps.SecretManager)
	useCases := credentialsUseCases.NewServices(credentialsUseCases.ServiceDependencies{
		TransactionManager: deps.TransactionManager,
		SecretManager:      deps.SecretManager,

		CredentialRepository:        repository,
		CredentialProcessor:         deps.NodeEngineCredentialProcessor,
		NodeEngineCredentialService: deps.NodeEngineCredentialService,
		PlatformCredentialService:   deps.PlatformCredentialService,
	})
	return &Module{
		CredentialService: service,
		useCases:          useCases,
	}
}

func (m *Module) Register(router fiber.Router, authMiddleware *auth.AuthMiddleware) {
	credentialsHTTP.RegisterRoutes(router, authMiddleware, m.useCases)
}
