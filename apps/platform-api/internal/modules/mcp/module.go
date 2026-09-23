package mcp

import (
	"slices"
	"time"

	"github.com/gofiber/fiber/v3"

	cachemiddleware "github.com/blocknextai/go-packages/fiber/middleware/cache"
	commonAuth "github.com/blocknextai/platform-api/internal/common/auth"
	commonSemaphore "github.com/blocknextai/platform-api/internal/common/semaphore"
	"github.com/blocknextai/platform-api/internal/config"
	apikeysContract "github.com/blocknextai/platform-api/internal/modules/apikeys/contract"
	credentialoauthContract "github.com/blocknextai/platform-api/internal/modules/credentialoauth/contract"
	credentialsContract "github.com/blocknextai/platform-api/internal/modules/credentials/contract"
	executionsContract "github.com/blocknextai/platform-api/internal/modules/executions/contract"
	"github.com/blocknextai/platform-api/internal/modules/mcp/internal/application/adapter"
	"github.com/blocknextai/platform-api/internal/modules/mcp/internal/application/credentialresolver"
	"github.com/blocknextai/platform-api/internal/modules/mcp/internal/application/history"
	"github.com/blocknextai/platform-api/internal/modules/mcp/internal/application/nodeservers"
	"github.com/blocknextai/platform-api/internal/modules/mcp/internal/application/servers"
	mcpHTTP "github.com/blocknextai/platform-api/internal/modules/mcp/internal/http"
	mcpUseCases "github.com/blocknextai/platform-api/internal/modules/mcp/internal/usecases"
	mcpoauthContract "github.com/blocknextai/platform-api/internal/modules/mcpoauth/contract"
	nodeengineContract "github.com/blocknextai/platform-api/internal/modules/nodeengine/contract"
)

type Dependencies struct {
	SemaphoreOptions config.SemaphoreOptions

	ServerURLTemplate string
	MaxExecutionTime  time.Duration

	ServerProviders []servers.ServerProvider

	ServerService                         nodeengineContract.ServerService
	ExecutorService                       nodeengineContract.ExecutorService
	CredentialService                     credentialsContract.CredentialService
	CredentialOAuthTokenRegenerateService credentialoauthContract.CredentialOAuthTokenRegenerateService
	ToolInvocationService                 executionsContract.ToolInvocationService
	MCPOAuthMetadataService               mcpoauthContract.MetadataService
}

type Module struct {
	useCases *mcpUseCases.Services
}

func NewModule(deps Dependencies) (*Module, error) {
	semaphoreManager, err := commonSemaphore.New(deps.SemaphoreOptions)
	if err != nil {
		return nil, err
	}

	resolver := credentialresolver.NewCredentialResolver(deps.CredentialService, deps.CredentialOAuthTokenRegenerateService)
	recorder := history.NewRecorder(deps.ToolInvocationService)
	mcpAdapter := adapter.NewAdapter(
		deps.ExecutorService,
		resolver,
		recorder,
		semaphoreManager,
		deps.SemaphoreOptions.MaxConcurrentExecutions,
		resolveHeartbeatInterval(deps.SemaphoreOptions),
		deps.MaxExecutionTime,
	)

	nodeServerProvider := nodeservers.NewServerProvider(deps.ServerService, mcpAdapter)

	useCases, err := mcpUseCases.NewServices(mcpUseCases.ServiceDependencies{
		ServerURLTemplate: deps.ServerURLTemplate,

		ServerProviders:         slices.Concat(deps.ServerProviders, []servers.ServerProvider{nodeServerProvider}),
		MCPOAuthMetadataService: deps.MCPOAuthMetadataService,
	})
	if err != nil {
		return nil, err
	}

	return &Module{
		useCases: useCases,
	}, nil
}

func (m *Module) Register(
	router fiber.Router,
	cacheMiddleware *cachemiddleware.Middleware,
	apiKeyMiddleware *commonAuth.APIKeyMiddleware[apikeysContract.Scope],
	accessTokenMiddleware *commonAuth.AccessTokenMiddleware,
) {
	mcpHTTP.RegisterRoutes(router, cacheMiddleware, m.useCases, apiKeyMiddleware, accessTokenMiddleware)
}

func resolveHeartbeatInterval(options config.SemaphoreOptions) time.Duration {
	if options.HeartbeatInterval > 0 {
		return options.HeartbeatInterval
	}

	return options.TTL / 3
}
