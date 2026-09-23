package clients

import (
	"github.com/blocknextai/go-packages/database"
	mcpOAuthDomainClients "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/domain/clients"
)

type Service struct {
	clientRepository   mcpOAuthDomainClients.ClientRepository
	transactionManager database.TransactionManager
}

func NewService(
	clientRepository mcpOAuthDomainClients.ClientRepository,
	transactionManager database.TransactionManager,
) *Service {
	return &Service{
		clientRepository:   clientRepository,
		transactionManager: transactionManager,
	}
}
