package usecases

import (
	"github.com/blocknextai/go-packages/database"
	"github.com/blocknextai/go-packages/secretmanager"
	credentialsDomainCredentials "github.com/blocknextai/platform-api/internal/modules/credentials/internal/domain/credentials"
	credentialsUseCases "github.com/blocknextai/platform-api/internal/modules/credentials/internal/usecases/credentials"
	nodeengineContract "github.com/blocknextai/platform-api/internal/modules/nodeengine/contract"
	platformContract "github.com/blocknextai/platform-api/internal/modules/platform/contract"
)

type Services struct {
	Credentials *credentialsUseCases.Service
}

type ServiceDependencies struct {
	TransactionManager database.TransactionManager
	SecretManager      secretmanager.SecretManager

	CredentialRepository        credentialsDomainCredentials.CredentialRepository
	CredentialProcessor         nodeengineContract.CredentialProcessor
	NodeEngineCredentialService nodeengineContract.CredentialService
	PlatformCredentialService   platformContract.PlatformCredentialService
}

func NewServices(deps ServiceDependencies) *Services {
	return &Services{
		Credentials: credentialsUseCases.NewService(deps.CredentialRepository, deps.SecretManager, deps.TransactionManager, deps.PlatformCredentialService, deps.CredentialProcessor, deps.NodeEngineCredentialService),
	}
}
