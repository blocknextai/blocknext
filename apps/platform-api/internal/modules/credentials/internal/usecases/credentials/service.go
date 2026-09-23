package credentials

import (
	"github.com/blocknextai/go-packages/database"
	"github.com/blocknextai/go-packages/secretmanager"
	credentialsDomainCredentials "github.com/blocknextai/platform-api/internal/modules/credentials/internal/domain/credentials"
	nodeengineContract "github.com/blocknextai/platform-api/internal/modules/nodeengine/contract"
	platformContract "github.com/blocknextai/platform-api/internal/modules/platform/contract"
)

type Service struct {
	credentialRepository        credentialsDomainCredentials.CredentialRepository
	secretManager               secretmanager.SecretManager
	transactionManager          database.TransactionManager
	platformCredentialService   platformContract.PlatformCredentialService
	credentialProcessor         nodeengineContract.CredentialProcessor
	nodeEngineCredentialService nodeengineContract.CredentialService
}

func NewService(
	credentialRepository credentialsDomainCredentials.CredentialRepository,
	secretManager secretmanager.SecretManager,
	transactionManager database.TransactionManager,
	platformCredentialService platformContract.PlatformCredentialService,
	credentialProcessor nodeengineContract.CredentialProcessor,
	nodeEngineCredentialService nodeengineContract.CredentialService,
) *Service {
	return &Service{
		credentialRepository:        credentialRepository,
		secretManager:               secretManager,
		transactionManager:          transactionManager,
		platformCredentialService:   platformCredentialService,
		credentialProcessor:         credentialProcessor,
		nodeEngineCredentialService: nodeEngineCredentialService,
	}
}
