package authorizationrequests

import (
	"time"

	"github.com/blocknextai/go-packages/database"
	mcpOAuthApplicationClients "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/application/clients"
	mcpOAuthDomainAuthorizationCodes "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/domain/authorizationcodes"
	mcpOAuthDomainAuthorizationRequests "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/domain/authorizationrequests"
	mcpOAuthDomainGrants "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/domain/grants"
	organizationsContract "github.com/blocknextai/platform-api/internal/modules/organizations/contract"
)

type Service struct {
	authorizationRequestRepository mcpOAuthDomainAuthorizationRequests.AuthorizationRequestRepository
	authorizationCodeRepository    mcpOAuthDomainAuthorizationCodes.AuthorizationCodeRepository
	grantRepository                mcpOAuthDomainGrants.GrantRepository
	transactionManager             database.TransactionManager
	organizationUserService        organizationsContract.OrganizationUserService
	issuerURL                      string
	authorizationCodeTTL           time.Duration
	clientResolver                 *mcpOAuthApplicationClients.ClientResolver
	resourceBaseURL                string
	consentURLTemplate             string
	authorizationRequestTTL        time.Duration
}

func NewService(
	authorizationRequestRepository mcpOAuthDomainAuthorizationRequests.AuthorizationRequestRepository,
	authorizationCodeRepository mcpOAuthDomainAuthorizationCodes.AuthorizationCodeRepository,
	grantRepository mcpOAuthDomainGrants.GrantRepository,
	transactionManager database.TransactionManager,
	organizationUserService organizationsContract.OrganizationUserService,
	issuerURL string,
	authorizationCodeTTL time.Duration,
	clientResolver *mcpOAuthApplicationClients.ClientResolver,
	resourceBaseURL string,
	consentURLTemplate string,
	authorizationRequestTTL time.Duration,
) *Service {
	return &Service{
		authorizationRequestRepository: authorizationRequestRepository,
		authorizationCodeRepository:    authorizationCodeRepository,
		grantRepository:                grantRepository,
		transactionManager:             transactionManager,
		organizationUserService:        organizationUserService,
		issuerURL:                      issuerURL,
		authorizationCodeTTL:           authorizationCodeTTL,
		clientResolver:                 clientResolver,
		resourceBaseURL:                resourceBaseURL,
		consentURLTemplate:             consentURLTemplate,
		authorizationRequestTTL:        authorizationRequestTTL,
	}
}
