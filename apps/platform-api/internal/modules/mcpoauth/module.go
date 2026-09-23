package mcpoauth

import (
	"database/sql"

	commonAuth "github.com/blocknextai/platform-api/internal/common/auth"

	"github.com/gofiber/fiber/v3"

	"github.com/blocknextai/go-packages/cache"
	"github.com/blocknextai/go-packages/database"
	cachemiddleware "github.com/blocknextai/go-packages/fiber/middleware/cache"
	"github.com/blocknextai/platform-api/internal/config"
	mcpOAuthApplicationClients "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/application/clients"
	mcpOAuthApplicationGrants "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/application/grants"
	mcpOAuthApplicationMetadata "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/application/metadata"
	mcpOAuthApplicationTokens "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/application/tokens"
	mcpOAuthDomainOAuth2 "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/domain/oauth2"
	mcpOAuthHTTP "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/http"
	mcpOAuthPostgresAuthorizationCodes "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/postgres/authorizationcodes"
	mcpOAuthPostgresAuthorizationRequests "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/postgres/authorizationrequests"
	mcpOAuthPostgresClients "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/postgres/clients"
	mcpOAuthPostgresGrants "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/postgres/grants"
	mcpOAuthPostgresRefreshTokens "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/postgres/refreshtokens"
	mcpOAuthTokenSigner "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/tokensigner"
	mcpOAuthUseCases "github.com/blocknextai/platform-api/internal/modules/mcpoauth/internal/usecases"
	organizationsContract "github.com/blocknextai/platform-api/internal/modules/organizations/contract"
)

type Dependencies struct {
	DB                 *sql.DB
	TransactionManager database.TransactionManager
	CacheService       cache.Service

	OAuthOptions config.MCPOAuthOptions

	OrganizationUserService organizationsContract.OrganizationUserService
}

type Module struct {
	AccessTokenValidator commonAuth.AccessTokenValidator
	MetadataService      mcpOAuthApplicationMetadata.MetadataService

	useCases *mcpOAuthUseCases.Services
}

func NewModule(deps Dependencies) (*Module, error) {
	if !mcpOAuthDomainOAuth2.IsValidBaseURL(deps.OAuthOptions.IssuerURL) || !mcpOAuthDomainOAuth2.IsValidBaseURL(deps.OAuthOptions.ResourceURL) {
		return nil, mcpOAuthDomainOAuth2.ErrInvalidBaseURL
	}

	tokenSigner, err := mcpOAuthTokenSigner.New(deps.OAuthOptions.IssuerURL, deps.OAuthOptions.SigningKey)
	if err != nil {
		return nil, err
	}

	clientRepository := mcpOAuthPostgresClients.NewClientRepository(deps.DB)
	authorizationRequestRepository := mcpOAuthPostgresAuthorizationRequests.NewAuthorizationRequestRepository(deps.DB)
	authorizationCodeRepository := mcpOAuthPostgresAuthorizationCodes.NewAuthorizationCodeRepository(deps.DB)
	grantRepository := mcpOAuthPostgresGrants.NewGrantRepository(deps.DB)
	refreshTokenRepository := mcpOAuthPostgresRefreshTokens.NewRefreshTokenRepository(deps.DB)
	metadataDocumentFetcher := mcpOAuthPostgresClients.NewMetadataDocumentFetcher(deps.CacheService)

	useCases := mcpOAuthUseCases.NewServices(mcpOAuthUseCases.ServiceDependencies{
		TransactionManager: deps.TransactionManager,

		OAuthOptions: deps.OAuthOptions,

		ClientRepository:               clientRepository,
		AuthorizationRequestRepository: authorizationRequestRepository,
		AuthorizationCodeRepository:    authorizationCodeRepository,
		GrantRepository:                grantRepository,
		RefreshTokenRepository:         refreshTokenRepository,
		ClientResolver:                 mcpOAuthApplicationClients.NewClientResolver(clientRepository, metadataDocumentFetcher),
		GrantRevoker:                   mcpOAuthApplicationGrants.NewGrantRevoker(grantRepository, refreshTokenRepository),
		AccessTokenSigner:              tokenSigner,
		AccessTokenValidator:           tokenSigner,
		OrganizationUserService:        deps.OrganizationUserService,
	})

	metadataService := mcpOAuthApplicationMetadata.NewMetadataService(deps.OAuthOptions.IssuerURL, deps.OAuthOptions.ResourceURL)

	return &Module{
		AccessTokenValidator: mcpOAuthApplicationTokens.NewResourceAccessTokenValidator(tokenSigner, metadataService),
		MetadataService:      metadataService,
		useCases:             useCases,
	}, nil
}

func (m *Module) Register(
	router fiber.Router,
	authMiddleware *commonAuth.AuthMiddleware,
	cacheMiddleware *cachemiddleware.Middleware,
) {
	mcpOAuthHTTP.RegisterRoutes(router, authMiddleware, cacheMiddleware, m.useCases)
}
