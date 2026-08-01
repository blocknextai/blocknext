package bootstrap

import (
	"context"
	"log/slog"

	"github.com/blocknextai/go-packages/auth/jwt"
	"github.com/blocknextai/platform-api/internal/account"
	"github.com/blocknextai/platform-api/internal/common"
	"github.com/blocknextai/platform-api/internal/config"
	"github.com/blocknextai/platform-api/internal/credentialoauth"
	"github.com/blocknextai/platform-api/internal/credentials"
	"github.com/blocknextai/platform-api/internal/executions"
	"github.com/blocknextai/platform-api/internal/llm"
	"github.com/blocknextai/platform-api/internal/nodeengine"
	"github.com/blocknextai/platform-api/internal/organizations"
	"github.com/blocknextai/platform-api/internal/platform"
	"github.com/blocknextai/platform-api/internal/realtime"
	"github.com/blocknextai/platform-api/internal/taskrunner"
	"github.com/blocknextai/platform-api/internal/triggers"
	"github.com/blocknextai/platform-api/internal/web3"
	"github.com/blocknextai/platform-api/internal/webhooks"
	"github.com/blocknextai/platform-api/internal/workflows"
)

type WebhookAPI struct {
	Core   *Core
	Config *config.WebhookAPIConfig

	Broadcaster           realtime.Broadcaster
	CommonModule          *common.Module
	Web3Module            *web3.Module
	AccountModule         *account.Module
	OrganizationsModule   *organizations.Module
	NodeEngineModule      *nodeengine.Module
	PlatformModule        *platform.Module
	CredentialsModule     *credentials.Module
	LLMModule             *llm.Module
	CredentialOAuthModule *credentialoauth.Module
	WorkflowsModule       *workflows.Module
	ExecutionsModule      *executions.Module
	TriggersModule        *triggers.Module
	TaskRunnerModule      *taskrunner.Module
	WebhooksModule        *webhooks.Module
}

func NewWebhookAPI(core *Core, cfg *config.WebhookAPIConfig) (*WebhookAPI, error) {
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

	llmModule, err := llm.NewModule(llm.Dependencies{
		ChatProviderOptions:            shared.Workflows.Generation,
		FunctionCallingProviderOptions: shared.FunctionCalling,
	})
	if err != nil {
		return nil, err
	}

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

	workflowsModule, err := workflows.NewModule(workflows.Dependencies{
		DB:                 core.DB,
		TransactionManager: core.TransactionManager,

		WorkflowsOptions: shared.Workflows,

		GenerationProvider:          llmModule.StreamingChatProvider,
		NodeEngineNodeService:       nodeEngineModule.NodeService,
		NodeEngineCredentialService: nodeEngineModule.CredentialService,
		NodeEngineAdapterService:    nodeEngineModule.AdapterService,
		OrganizationUserService:     organizationsModule.OrganizationUserService,
		UserService:                 accountModule.UserService,
		LinkedAccountService:        accountModule.LinkedAccountService,
	})
	if err != nil {
		return nil, err
	}

	executionsModule := executions.NewModule(executions.Dependencies{
		DB:                 core.DB,
		TransactionManager: core.TransactionManager,

		OrganizationUserService: organizationsModule.OrganizationUserService,
		WorkflowService:         workflowsModule.WorkflowService,
		UserService:             accountModule.UserService,
		LinkedAccountService:    accountModule.LinkedAccountService,
	})

	triggersModule := triggers.NewModule(triggers.Dependencies{
		DB:                 core.DB,
		TransactionManager: core.TransactionManager,
		SecretManager:      core.SecretManager,

		WorkflowService: workflowsModule.WorkflowService,
	})

	taskRunnerModule, err := taskrunner.NewModule(taskrunner.Dependencies{
		CacheService: core.CacheService,
		Broadcaster:  broadcaster,

		TaskRunnerOptions: cfg.TaskRunner,

		FunctionCallingService:                llmModule.FunctionCallingService,
		CredentialOAuthTokenRegenerateService: credentialOAuthModule.CredentialOAuthTokenRegenerateService,
		WorkflowService:                       workflowsModule.WorkflowService,
		TaskExecutionService:                  executionsModule.TaskExecutionService,
		TaskClaimService:                      executionsModule.TaskClaimService,
		NodeExecutionService:                  executionsModule.NodeExecutionService,
		TriggerService:                        triggersModule.TriggerService,
		WebhookResolver:                       triggersModule.WebhookResolver,
	})
	if err != nil {
		return nil, err
	}

	webhooksModule := webhooks.NewModule(webhooks.Dependencies{
		TriggerWebhookProcessor: taskRunnerModule.WebhookProcessor,
	})

	return &WebhookAPI{
		Core:                  core,
		Config:                cfg,
		Broadcaster:           broadcaster,
		CommonModule:          commonModule,
		Web3Module:            web3Module,
		AccountModule:         accountModule,
		OrganizationsModule:   organizationsModule,
		NodeEngineModule:      nodeEngineModule,
		PlatformModule:        platformModule,
		CredentialsModule:     credentialsModule,
		LLMModule:             llmModule,
		CredentialOAuthModule: credentialOAuthModule,
		WorkflowsModule:       workflowsModule,
		ExecutionsModule:      executionsModule,
		TriggersModule:        triggersModule,
		TaskRunnerModule:      taskRunnerModule,
		WebhooksModule:        webhooksModule,
	}, nil
}

func (w *WebhookAPI) Health(ctx context.Context) error {
	if err := w.Core.Health(ctx); err != nil {
		return err
	}
	if err := w.Broadcaster.Ping(ctx); err != nil {
		return err
	}
	return w.TaskRunnerModule.Health(ctx)
}

func (w *WebhookAPI) Shutdown() error {
	if err := w.TaskRunnerModule.Shutdown(); err != nil {
		slog.Error("failed to shutdown task runner module", "error", err)
		return err
	}
	if err := w.Broadcaster.Close(); err != nil {
		slog.Error("failed to close broadcaster", "error", err)
		return err
	}
	return nil
}
