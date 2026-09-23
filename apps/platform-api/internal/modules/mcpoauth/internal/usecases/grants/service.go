package grants

import (
	"github.com/blocknextai/go-packages/database"
	mcpOAuthApplicationClients "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/application/clients"
	mcpOAuthApplicationGrants "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/application/grants"
	mcpOAuthDomainGrants "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/domain/grants"
)

type Service struct {
	grantRepository    mcpOAuthDomainGrants.GrantRepository
	clientResolver     *mcpOAuthApplicationClients.ClientResolver
	grantRevoker       *mcpOAuthApplicationGrants.GrantRevoker
	transactionManager database.TransactionManager
}

func NewService(
	grantRepository mcpOAuthDomainGrants.GrantRepository,
	clientResolver *mcpOAuthApplicationClients.ClientResolver,
	grantRevoker *mcpOAuthApplicationGrants.GrantRevoker,
	transactionManager database.TransactionManager,
) *Service {
	return &Service{
		grantRepository:    grantRepository,
		clientResolver:     clientResolver,
		grantRevoker:       grantRevoker,
		transactionManager: transactionManager,
	}
}
