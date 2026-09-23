package apikeys

import (
	"github.com/blocknextai/go-packages/database"
	apiKeysDomainAPIKeys "github.com/blocknextai/platform-api/internal/modules/apikeys/internal/domain/apikeys"
)

type Service struct {
	apiKeyRepository   apiKeysDomainAPIKeys.APIKeyRepository
	transactionManager database.TransactionManager
}

func NewService(
	apiKeyRepository apiKeysDomainAPIKeys.APIKeyRepository,
	transactionManager database.TransactionManager,
) *Service {
	return &Service{
		apiKeyRepository:   apiKeyRepository,
		transactionManager: transactionManager,
	}
}
