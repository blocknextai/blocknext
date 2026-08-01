package credentials

import (
	"database/sql"

	"github.com/blocknextai/go-packages/database"
	"github.com/blocknextai/go-packages/secretmanager"
	"github.com/blocknextai/platform-api/internal/common/presentation/auth"
	credentialsApplicationCredentials "github.com/blocknextai/platform-api/internal/credentials/application/credentials"
	credentialsInfrastructure "github.com/blocknextai/platform-api/internal/credentials/infrastructure"
	credentialsInfrastructureCredentials "github.com/blocknextai/platform-api/internal/credentials/infrastructure/credentials"
	credentialsPresentation "github.com/blocknextai/platform-api/internal/credentials/presentation"
	nodeEngineApplicationCredentials "github.com/blocknextai/platform-api/internal/nodeengine/application/credentials"
	platformApplicationPlatformCredentials "github.com/blocknextai/platform-api/internal/platform/application/platformcredentials"
	"github.com/gofiber/fiber/v3"
)

type Dependencies struct {
	DB                 *sql.DB
	TransactionManager database.TransactionManager
	SecretManager      secretmanager.SecretManager

	NodeEngineCredentialProcessor nodeEngineApplicationCredentials.CredentialProcessor
	NodeEngineCredentialService   nodeEngineApplicationCredentials.CredentialService
	PlatformCredentialService     platformApplicationPlatformCredentials.PlatformCredentialService
}

type Module struct {
	CredentialService credentialsApplicationCredentials.CredentialService

	handlers *credentialsInfrastructure.Handlers
}

func NewModule(deps Dependencies) *Module {
	repository := credentialsInfrastructureCredentials.NewCredentialRepository(deps.DB)
	service := credentialsApplicationCredentials.NewCredentialService(repository, deps.SecretManager)
	handlers := credentialsInfrastructure.RegisterInfrastructure(credentialsInfrastructure.RegisterInfrastructureDeps{
		TransactionManager: deps.TransactionManager,
		SecretManager:      deps.SecretManager,

		CredentialRepository:        repository,
		CredentialProcessor:         deps.NodeEngineCredentialProcessor,
		NodeEngineCredentialService: deps.NodeEngineCredentialService,
		PlatformCredentialService:   deps.PlatformCredentialService,
	})
	return &Module{
		CredentialService: service,
		handlers:          handlers,
	}
}

func (m *Module) Register(router fiber.Router, authMiddleware *auth.AuthMiddleware) {
	credentialsPresentation.RegisterPresentation(router, authMiddleware, m.handlers)
}
