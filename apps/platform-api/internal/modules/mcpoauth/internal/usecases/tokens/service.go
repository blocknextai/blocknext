package tokens

import (
	"time"

	"github.com/blocknextai/go-packages/database"
	mcpOAuthApplicationClients "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/application/clients"
	mcpOAuthApplicationGrants "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/application/grants"
	mcpOAuthApplicationTokens "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/application/tokens"
	mcpOAuthDomainAuthorizationCodes "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/domain/authorizationcodes"
	mcpOAuthDomainGrants "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/domain/grants"
	mcpOAuthDomainRefreshTokens "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/domain/refreshtokens"
)

type Service struct {
	authorizationCodeRepository mcpOAuthDomainAuthorizationCodes.AuthorizationCodeRepository
	refreshTokenRepository      mcpOAuthDomainRefreshTokens.RefreshTokenRepository
	grantRepository             mcpOAuthDomainGrants.GrantRepository
	clientResolver              *mcpOAuthApplicationClients.ClientResolver
	grantRevoker                *mcpOAuthApplicationGrants.GrantRevoker
	accessTokenSigner           mcpOAuthApplicationTokens.AccessTokenSigner
	transactionManager          database.TransactionManager
	resourceBaseURL             string
	accessTokenTTL              time.Duration
	refreshTokenTTL             time.Duration
	accessTokenValidator        mcpOAuthApplicationTokens.AccessTokenValidator
}

func NewService(
	authorizationCodeRepository mcpOAuthDomainAuthorizationCodes.AuthorizationCodeRepository,
	refreshTokenRepository mcpOAuthDomainRefreshTokens.RefreshTokenRepository,
	grantRepository mcpOAuthDomainGrants.GrantRepository,
	clientResolver *mcpOAuthApplicationClients.ClientResolver,
	grantRevoker *mcpOAuthApplicationGrants.GrantRevoker,
	accessTokenSigner mcpOAuthApplicationTokens.AccessTokenSigner,
	transactionManager database.TransactionManager,
	resourceBaseURL string,
	accessTokenTTL time.Duration,
	refreshTokenTTL time.Duration,
	accessTokenValidator mcpOAuthApplicationTokens.AccessTokenValidator,
) *Service {
	return &Service{
		authorizationCodeRepository: authorizationCodeRepository,
		refreshTokenRepository:      refreshTokenRepository,
		grantRepository:             grantRepository,
		clientResolver:              clientResolver,
		grantRevoker:                grantRevoker,
		accessTokenSigner:           accessTokenSigner,
		transactionManager:          transactionManager,
		resourceBaseURL:             resourceBaseURL,
		accessTokenTTL:              accessTokenTTL,
		refreshTokenTTL:             refreshTokenTTL,
		accessTokenValidator:        accessTokenValidator,
	}
}
