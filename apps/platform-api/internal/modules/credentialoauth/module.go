package credentialoauth

import (
	"time"

	"github.com/gofiber/fiber/v3"

	"github.com/blocknextai/go-packages/cache"
	"github.com/blocknextai/platform-api/internal/common/auth"
	credentialOAuthApplicationOAuth2 "github.com/blocknextai/platform-api/internal/modules/credentialoauth/internal/application/oauth2"
	credentialOAuthApplicationRegenerate "github.com/blocknextai/platform-api/internal/modules/credentialoauth/internal/application/regenerate"
	credentialOAuthHTTP "github.com/blocknextai/platform-api/internal/modules/credentialoauth/internal/http"
	credentialOAuthUseCases "github.com/blocknextai/platform-api/internal/modules/credentialoauth/internal/usecases"
	credentialsContract "github.com/blocknextai/platform-api/internal/modules/credentials/contract"
	nodeengineContract "github.com/blocknextai/platform-api/internal/modules/nodeengine/contract"
	platformContract "github.com/blocknextai/platform-api/internal/modules/platform/contract"
)

type Dependencies struct {
	CacheService cache.Service

	StateTTL          time.Duration
	OAuth2RedirectURL string

	CredentialService           credentialsContract.CredentialService
	NodeEngineCredentialService nodeengineContract.CredentialService
	PlatformCredentialService   platformContract.PlatformCredentialService
}

type Module struct {
	CredentialOAuthTokenRegenerateService credentialOAuthApplicationRegenerate.CredentialOAuthTokenRegenerateService

	useCases *credentialOAuthUseCases.Services
}

func NewModule(deps Dependencies) *Module {
	stateStore := credentialOAuthApplicationOAuth2.NewStateStore(deps.CacheService, deps.StateTTL)

	regenerateService := credentialOAuthApplicationRegenerate.NewCredentialOAuthTokenRegenerateService(
		deps.CredentialService,
		deps.NodeEngineCredentialService,
		deps.PlatformCredentialService,
		deps.CacheService,
	)

	useCases := credentialOAuthUseCases.NewServices(credentialOAuthUseCases.ServiceDependencies{
		OAuth2RedirectURL: deps.OAuth2RedirectURL,

		CredentialService:           deps.CredentialService,
		NodeEngineCredentialService: deps.NodeEngineCredentialService,
		PlatformCredentialService:   deps.PlatformCredentialService,
		StateStore:                  stateStore,
	})
	return &Module{
		CredentialOAuthTokenRegenerateService: regenerateService,
		useCases:                              useCases,
	}
}

func (m *Module) Register(router fiber.Router, authMiddleware *auth.AuthMiddleware) {
	credentialOAuthHTTP.RegisterRoutes(router, authMiddleware, m.useCases)
}
