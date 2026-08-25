package mcp

import (
	cachemiddleware "github.com/blocknextai/go-packages/fiber/middleware/cache"
	commonInfrastructureSemaphore "github.com/blocknextai/platform-api/internal/common/infrastructure/semaphore"
	commonPresentationAuth "github.com/blocknextai/platform-api/internal/common/presentation/auth"
	"github.com/blocknextai/platform-api/internal/config"
	credentialOAuthApplicationRegenerate "github.com/blocknextai/platform-api/internal/credentialoauth/application/regenerate"
	executionsApplicationToolInvocations "github.com/blocknextai/platform-api/internal/executions/application/toolinvocations"
	"github.com/blocknextai/platform-api/internal/mcp/application/adapter"
	"github.com/blocknextai/platform-api/internal/mcp/application/credentialresolver"
	"github.com/blocknextai/platform-api/internal/mcp/application/history"
	mcpInfrastructure "github.com/blocknextai/platform-api/internal/mcp/infrastructure"
	mcpPresentation "github.com/blocknextai/platform-api/internal/mcp/presentation"
	nodeEngineApplicationExecutors "github.com/blocknextai/platform-api/internal/nodeengine/application/executors"
	nodeEngineApplicationMCP "github.com/blocknextai/platform-api/internal/nodeengine/application/mcp"
	"github.com/blocknextai/platform-api/internal/realtime"
	"github.com/gofiber/fiber/v3"
	"time"
)

type Dependencies struct {
	ServerURLTemplate string
	MaxExecutionTime  time.Duration

	SemaphoreOptions config.SemaphoreOptions

	ServerService                         nodeEngineApplicationMCP.ServerService
	ExecutorService                       nodeEngineApplicationExecutors.ExecutorService
	CredentialOAuthTokenRegenerateService credentialOAuthApplicationRegenerate.CredentialOAuthTokenRegenerateService
	ToolInvocationService                 executionsApplicationToolInvocations.ToolInvocationService
	Broadcaster                           realtime.Broadcaster
}

type Module struct {
	handlers *mcpInfrastructure.Handlers
}

func NewModule(deps Dependencies) (*Module, error) {
	semaphoreManager, err := commonInfrastructureSemaphore.New(deps.SemaphoreOptions)
	if err != nil {
		return nil, err
	}

	resolver := credentialresolver.NewCredentialResolver(deps.CredentialOAuthTokenRegenerateService)
	recorder := history.NewRecorder(deps.ToolInvocationService, deps.Broadcaster)
	mcpAdapter := adapter.NewAdapter(
		deps.ExecutorService,
		resolver,
		recorder,
		semaphoreManager,
		deps.SemaphoreOptions.MaxConcurrentExecutions,
		resolveHeartbeatInterval(deps.SemaphoreOptions),
		deps.MaxExecutionTime,
	)

	handlers, err := mcpInfrastructure.RegisterInfrastructure(mcpInfrastructure.RegisterInfrastructureDeps{
		ServerURLTemplate: deps.ServerURLTemplate,

		ServerService: deps.ServerService,
		Adapter:       mcpAdapter,
	})
	if err != nil {
		return nil, err
	}

	return &Module{
		handlers: handlers,
	}, nil
}

func (m *Module) Register(
	router fiber.Router,
	cacheMiddleware *cachemiddleware.Middleware,
	apiKeyMiddleware *commonPresentationAuth.APIKeyMiddleware,
) {
	mcpPresentation.RegisterPresentation(router, cacheMiddleware, m.handlers, apiKeyMiddleware)
}

func resolveHeartbeatInterval(options config.SemaphoreOptions) time.Duration {
	if options.HeartbeatInterval > 0 {
		return options.HeartbeatInterval
	}

	return options.TTL / 3
}
