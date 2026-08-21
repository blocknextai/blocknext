package bootstrap

import (
	"context"

	"github.com/blocknextai/go-packages/auth/jwt"
	"github.com/blocknextai/platform-api/internal/account"
	"github.com/blocknextai/platform-api/internal/apikeys"
	"github.com/blocknextai/platform-api/internal/common"
	"github.com/blocknextai/platform-api/internal/config"
	"github.com/blocknextai/platform-api/internal/credentialoauth"
	"github.com/blocknextai/platform-api/internal/credentials"
	"github.com/blocknextai/platform-api/internal/mcp"
	"github.com/blocknextai/platform-api/internal/nodeengine"
	"github.com/blocknextai/platform-api/internal/organizations"
	"github.com/blocknextai/platform-api/internal/platform"
	"github.com/blocknextai/platform-api/internal/web3"
)

type MCPAPI struct {
	Core   *Core
	Config *config.MCPAPIConfig

	CommonModule          *common.Module
	Web3Module            *web3.Module
	AccountModule         *account.Module
	OrganizationsModule   *organizations.Module
	APIKeysModule         *apikeys.Module
	NodeEngineModule      *nodeengine.Module
	PlatformModule        *platform.Module
	CredentialsModule     *credentials.Module
	CredentialOAuthModule *credentialoauth.Module
	MCPModule             *mcp.Module
}

func NewMCPAPI(core *Core, cfg *config.MCPAPIConfig) (*MCPAPI, error) {
	shared := cfg.Shared

	jwtService, err := jwt.New(
		shared.JWT.Issuer,
		shared.JWT.Audience,
		shared.JWT.SecretKey,
		shared.JWT.AccessTokenExpirationTime,
		shared.JWT.RefreshTokenExpirationTime,
		shared.JWT.Leeway,
	)
	if err != nil {
		return nil, err
	}

	commonModule := common.NewModule(common.Dependencies{
		EmailSenderOptions: shared.EmailSender,
		BcryptCost:         shared.Auth.Password.BcryptCost,
	})

	web3Module, err := web3.NewModule(web3.Dependencies{
		LoginMessage: shared.Auth.Metamask.LoginMessage,
	})
	if err != nil {
		return nil, err
	}

	accountModule, err := account.NewModule(account.Dependencies{
		DB:                       core.DB,
		TransactionManager:       core.TransactionManager,
		EventBus:                 core.EventBus.Bus,
		EventBusPublisherService: core.EventBus.PublisherService,
		EventBusInboxService:     core.EventBus.InboxService,
		CacheService:             core.CacheService,

		AccessTokenExpirationTime: shared.JWT.AccessTokenExpirationTime,
		AuthOptions:               shared.Auth,
		PlatformUIBaseURL:         shared.PlatformUI.BaseURL,

		JWTService:        jwtService,
		EmailSender:       commonModule.EmailSender,
		PasswordHasher:    commonModule.PasswordHasher,
		SignatureVerifier: web3Module.SignatureVerifier,
	})
	if err != nil {
		return nil, err
	}

	nodeEngineModule := nodeengine.NewModule(nodeengine.Dependencies{
		FileGateway: core.FileGateway,

		OAuth2RedirectURL:         shared.CredentialOAuth.OAuth2RedirectURL,
		WebhookTriggerURLTemplate: shared.Webhook.Trigger.URLTemplate,
	})

	platformModule := platform.NewModule(platform.Dependencies{
		CredentialConfigs: shared.Platform.Credentials.ToCredentialConfigs(),

		NodeEngineCredentialService: nodeEngineModule.CredentialService,
	})

	credentialsModule := credentials.NewModule(credentials.Dependencies{
		DB:                 core.DB,
		TransactionManager: core.TransactionManager,
		SecretManager:      core.SecretManager,

		NodeEngineCredentialProcessor: nodeEngineModule.CredentialProcessor,
		NodeEngineCredentialService:   nodeEngineModule.CredentialService,
		PlatformCredentialService:     platformModule.PlatformCredentialService,
	})

	credentialOAuthModule := credentialoauth.NewModule(credentialoauth.Dependencies{
		CacheService: core.CacheService,

		StateTTL:          shared.CredentialOAuth.StateTTL,
		OAuth2RedirectURL: shared.CredentialOAuth.OAuth2RedirectURL,

		CredentialService:           credentialsModule.CredentialService,
		NodeEngineCredentialService: nodeEngineModule.CredentialService,
		PlatformCredentialService:   platformModule.PlatformCredentialService,
	})

	organizationsModule := organizations.NewModule(organizations.Dependencies{
		DB:                       core.DB,
		TransactionManager:       core.TransactionManager,
		EventBusPublisherService: core.EventBus.PublisherService,
		CacheService:             core.CacheService,

		UserService:          accountModule.UserService,
		LinkedAccountService: accountModule.LinkedAccountService,
	})

	apiKeysModule := apikeys.NewModule(apikeys.Dependencies{
		DB:                 core.DB,
		TransactionManager: core.TransactionManager,
	})

	mcpModule, err := mcp.NewModule(mcp.Dependencies{
		ServerURLTemplate: cfg.MCP.Server.URLTemplate,

		ServerService:                         nodeEngineModule.MCPServerService,
		ExecutorService:                       nodeEngineModule.ExecutorService,
		CredentialOAuthTokenRegenerateService: credentialOAuthModule.CredentialOAuthTokenRegenerateService,
	})
	if err != nil {
		return nil, err
	}

	return &MCPAPI{
		Core:                  core,
		Config:                cfg,
		CommonModule:          commonModule,
		Web3Module:            web3Module,
		AccountModule:         accountModule,
		OrganizationsModule:   organizationsModule,
		APIKeysModule:         apiKeysModule,
		NodeEngineModule:      nodeEngineModule,
		PlatformModule:        platformModule,
		CredentialsModule:     credentialsModule,
		CredentialOAuthModule: credentialOAuthModule,
		MCPModule:             mcpModule,
	}, nil
}

func (a *MCPAPI) Health(ctx context.Context) error {
	return a.Core.Health(ctx)
}
