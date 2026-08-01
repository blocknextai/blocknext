package credentialoauth

import (
	"time"

	"github.com/blocknextai/go-packages/cache"
	"github.com/blocknextai/platform-api/internal/common/presentation/auth"
	credentialOAuthApplicationOAuth2 "github.com/blocknextai/platform-api/internal/credentialoauth/application/oauth2"
	credentialOAuthApplicationRegenerate "github.com/blocknextai/platform-api/internal/credentialoauth/application/regenerate"
	credentialOAuthInfrastructure "github.com/blocknextai/platform-api/internal/credentialoauth/infrastructure"
	credentialOAuthPresentation "github.com/blocknextai/platform-api/internal/credentialoauth/presentation"
	credentialsApplicationCredentials "github.com/blocknextai/platform-api/internal/credentials/application/credentials"
	nodeEngineApplicationCredentials "github.com/blocknextai/platform-api/internal/nodeengine/application/credentials"
	platformApplicationPlatformCredentials "github.com/blocknextai/platform-api/internal/platform/application/platformcredentials"
	"github.com/gofiber/fiber/v3"
)

type Dependencies struct {
	CacheService cache.Service

	StateTTL          time.Duration
	OAuth2RedirectURL string

	CredentialService           credentialsApplicationCredentials.CredentialService
	NodeEngineCredentialService nodeEngineApplicationCredentials.CredentialService
	PlatformCredentialService   platformApplicationPlatformCredentials.PlatformCredentialService
}

type Module struct {
	CredentialOAuthTokenRegenerateService credentialOAuthApplicationRegenerate.CredentialOAuthTokenRegenerateService

	handlers *credentialOAuthInfrastructure.Handlers
}

func NewModule(deps Dependencies) *Module {
	stateStore := credentialOAuthApplicationOAuth2.NewStateStore(deps.CacheService, deps.StateTTL)

	regenerateService := credentialOAuthApplicationRegenerate.NewCredentialOAuthTokenRegenerateService(
		deps.CredentialService,
		deps.NodeEngineCredentialService,
		deps.PlatformCredentialService,
		deps.CacheService,
	)

	handlers := credentialOAuthInfrastructure.RegisterInfrastructure(credentialOAuthInfrastructure.RegisterInfrastructureDeps{
		OAuth2RedirectURL: deps.OAuth2RedirectURL,

		CredentialService:           deps.CredentialService,
		NodeEngineCredentialService: deps.NodeEngineCredentialService,
		PlatformCredentialService:   deps.PlatformCredentialService,
		StateStore:                  stateStore,
	})
	return &Module{
		CredentialOAuthTokenRegenerateService: regenerateService,
		handlers:                              handlers,
	}
}

func (m *Module) Register(router fiber.Router, authMiddleware *auth.AuthMiddleware) {
	credentialOAuthPresentation.RegisterPresentation(router, authMiddleware, m.handlers)
}
