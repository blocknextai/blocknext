package usecases

import (
	"github.com/blocknextai/go-packages/database"
	apiKeysDomainAPIKeys "github.com/blocknextai/platform-api/internal/modules/apikeys/internal/domain/apikeys"
	apikeysUseCases "github.com/blocknextai/platform-api/internal/modules/apikeys/internal/usecases/apikeys"
)

type Services struct {
	APIKeys *apikeysUseCases.Service
}

type ServiceDependencies struct {
	TransactionManager database.TransactionManager

	ApiKeyRepository apiKeysDomainAPIKeys.APIKeyRepository
}

func NewServices(deps ServiceDependencies) *Services {
	return &Services{
		APIKeys: apikeysUseCases.NewService(deps.ApiKeyRepository, deps.TransactionManager),
	}
}
