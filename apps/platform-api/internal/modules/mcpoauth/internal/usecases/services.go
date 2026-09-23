package usecases

import (
	"github.com/blocknextai/go-packages/database"
	"github.com/blocknextai/platform-api/internal/config"
	mcpOAuthApplicationClients "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/application/clients"
	mcpOAuthApplicationGrants "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/application/grants"
	mcpOAuthApplicationTokens "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/application/tokens"
	mcpOAuthDomainAuthorizationCodes "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/domain/authorizationcodes"
	mcpOAuthDomainAuthorizationRequests "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/domain/authorizationrequests"
	mcpOAuthDomainClients "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/domain/clients"
	mcpOAuthDomainGrants "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/domain/grants"
	mcpOAuthDomainRefreshTokens "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/domain/refreshtokens"
	authorizationrequestsUseCases "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/usecases/authorizationrequests"
	clientsUseCases "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/usecases/clients"
	grantsUseCases "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/usecases/grants"
	metadataUseCases "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/usecases/metadata"
	tokensUseCases "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/usecases/tokens"
	organizationsContract "github.com/blocknextai/platform-api/internal/modules/organizations/contract"
)

type Services struct {
	AuthorizationRequests *authorizationrequestsUseCases.Service
	Clients               *clientsUseCases.Service
	Grants                *grantsUseCases.Service
	Metadata              *metadataUseCases.Service
	Tokens                *tokensUseCases.Service
}

type ServiceDependencies struct {
	TransactionManager database.TransactionManager

	OAuthOptions config.MCPOAuthOptions

	ClientRepository               mcpOAuthDomainClients.ClientRepository
	AuthorizationRequestRepository mcpOAuthDomainAuthorizationRequests.AuthorizationRequestRepository
	AuthorizationCodeRepository    mcpOAuthDomainAuthorizationCodes.AuthorizationCodeRepository
	GrantRepository                mcpOAuthDomainGrants.GrantRepository
	RefreshTokenRepository         mcpOAuthDomainRefreshTokens.RefreshTokenRepository
	ClientResolver                 *mcpOAuthApplicationClients.ClientResolver
	GrantRevoker                   *mcpOAuthApplicationGrants.GrantRevoker
	AccessTokenSigner              mcpOAuthApplicationTokens.AccessTokenSigner
	AccessTokenValidator           mcpOAuthApplicationTokens.AccessTokenValidator
	OrganizationUserService        organizationsContract.OrganizationUserService
}

func NewServices(deps ServiceDependencies) *Services {
	options := deps.OAuthOptions

	return &Services{
		AuthorizationRequests: authorizationrequestsUseCases.NewService(deps.AuthorizationRequestRepository, deps.AuthorizationCodeRepository, deps.GrantRepository, deps.TransactionManager, deps.OrganizationUserService, options.IssuerURL, options.AuthorizationCodeTTL, deps.ClientResolver, options.ResourceURL, options.ConsentURLTemplate, options.AuthorizationRequestTTL),
		Clients:               clientsUseCases.NewService(deps.ClientRepository, deps.TransactionManager),
		Grants:                grantsUseCases.NewService(deps.GrantRepository, deps.ClientResolver, deps.GrantRevoker, deps.TransactionManager),
		Metadata:              metadataUseCases.NewService(options.IssuerURL),
		Tokens:                tokensUseCases.NewService(deps.AuthorizationCodeRepository, deps.RefreshTokenRepository, deps.GrantRepository, deps.ClientResolver, deps.GrantRevoker, deps.AccessTokenSigner, deps.TransactionManager, options.ResourceURL, options.AccessTokenTTL, options.RefreshTokenTTL, deps.AccessTokenValidator),
	}
}
