package bootstrap

import (
	"context"

	"github.com/blocknextai/go-packages/auth/jwt"
	"github.com/blocknextai/platform-api/internal/common"
	"github.com/blocknextai/platform-api/internal/config"
	"github.com/blocknextai/platform-api/internal/modules/account"
	"github.com/blocknextai/platform-api/internal/modules/apikeys"
	"github.com/blocknextai/platform-api/internal/modules/credentialoauth"
	"github.com/blocknextai/platform-api/internal/modules/credentials"
	"github.com/blocknextai/platform-api/internal/modules/executions"
	"github.com/blocknextai/platform-api/internal/modules/mcp"
	mcpContract "github.com/blocknextai/platform-api/internal/modules/mcp/contract"
	"github.com/blocknextai/platform-api/internal/modules/mcpoauth"
	"github.com/blocknextai/platform-api/internal/modules/mcpplatform"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine"
	"github.com/blocknextai/platform-api/internal/modules/organizations"
	"github.com/blocknextai/platform-api/internal/modules/platform"
	"github.com/blocknextai/platform-api/internal/modules/triggers"
	"github.com/blocknextai/platform-api/internal/modules/workflows"
	"github.com/blocknextai/platform-api/internal/realtime"
)

type MCPAPI struct {
	Core   *Core
	Config *config.MCPAPIConfig

	JWTService jwt.AuthJWTService

	CommonModule          *common.Module
	AccountModule         *account.Module
	OrganizationsModule   *organizations.Module
	APIKeysModule         *apikeys.Module
	NodeEngineModule      *nodeengine.Module
	PlatformModule        *platform.Module
	CredentialsModule     *credentials.Module
	CredentialOAuthModule *credentialoauth.Module
	MCPOAuthModule        *mcpoauth.Module
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

	broadcaster, err := realtime.New(shared.Broker)
	if err != nil {
		return nil, err
	}

	commonModule := common.NewModule(common.Dependencies{
		EmailSenderOptions: shared.EmailSender,
		BcryptCost:         shared.Auth.Password.BcryptCost,
	})

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

		JWTService:     jwtService,
		EmailSender:    commonModule.EmailSender,
		PasswordHasher: commonModule.PasswordHasher,
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

	executionServices := executions.NewServices(executions.ServicesDependencies{
		DB:                 core.DB,
		TransactionManager: core.TransactionManager,
		Broadcaster:        broadcaster,

		OrganizationUserService: organizationsModule.OrganizationUserService,
	})

	mcpOAuthModule, err := mcpoauth.NewModule(mcpoauth.Dependencies{
		DB:                 core.DB,
		TransactionManager: core.TransactionManager,
		CacheService:       core.CacheService,

		OAuthOptions: cfg.MCP.OAuth,

		OrganizationUserService: organizationsModule.OrganizationUserService,
	})
	if err != nil {
		return nil, err
	}

	workflowServices := workflows.NewServices(workflows.ServicesDependencies{
		DB: core.DB,
	})

	triggerServices := triggers.NewServices(triggers.ServicesDependencies{
		DB: core.DB,
	})

	mcpPlatformModule := mcpplatform.NewModule(mcpplatform.Dependencies{
		UserService:             accountModule.UserService,
		OrganizationService:     organizationsModule.OrganizationService,
		OrganizationUserService: organizationsModule.OrganizationUserService,
		WorkflowService:         workflowServices.WorkflowService,
		TriggerService:          triggerServices.TriggerService,
		TaskExecutionService:    executionServices.TaskExecutionService,
		ToolInvocationService:   executionServices.ToolInvocationService,
		CredentialService:       credentialsModule.CredentialService,
	})

	mcpModule, err := mcp.NewModule(mcp.Dependencies{
		SemaphoreOptions: shared.Semaphore,

		ServerURLTemplate: cfg.MCP.Server.URLTemplate,
		MaxExecutionTime:  cfg.MCP.MaxExecutionTime,

		ServerProviders: []mcpContract.ServerProvider{
			mcpPlatformModule.ServerProvider,
		},

		ServerService:                         nodeEngineModule.MCPServerService,
		ExecutorService:                       nodeEngineModule.ExecutorService,
		CredentialService:                     credentialsModule.CredentialService,
		CredentialOAuthTokenRegenerateService: credentialOAuthModule.CredentialOAuthTokenRegenerateService,
		ToolInvocationService:                 executionServices.ToolInvocationService,
		MCPOAuthMetadataService:               mcpOAuthModule.MetadataService,
	})
	if err != nil {
		return nil, err
	}

	return &MCPAPI{
		Core:                  core,
		Config:                cfg,
		JWTService:            jwtService,
		CommonModule:          commonModule,
		AccountModule:         accountModule,
		OrganizationsModule:   organizationsModule,
		APIKeysModule:         apiKeysModule,
		NodeEngineModule:      nodeEngineModule,
		PlatformModule:        platformModule,
		CredentialsModule:     credentialsModule,
		CredentialOAuthModule: credentialOAuthModule,
		MCPOAuthModule:        mcpOAuthModule,
		MCPModule:             mcpModule,
	}, nil
}

func (a *MCPAPI) Health(ctx context.Context) error {
	return a.Core.Health(ctx)
}
