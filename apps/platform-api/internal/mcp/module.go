package mcp

import (
	cachemiddleware "github.com/blocknextai/go-packages/fiber/middleware/cache"
	commonPresentationAuth "github.com/blocknextai/platform-api/internal/common/presentation/auth"
	credentialOAuthApplicationRegenerate "github.com/blocknextai/platform-api/internal/credentialoauth/application/regenerate"
	"github.com/blocknextai/platform-api/internal/mcp/application/adapter"
	"github.com/blocknextai/platform-api/internal/mcp/application/credentialresolver"
	mcpInfrastructure "github.com/blocknextai/platform-api/internal/mcp/infrastructure"
	mcpPresentation "github.com/blocknextai/platform-api/internal/mcp/presentation"
	nodeEngineApplicationExecutors "github.com/blocknextai/platform-api/internal/nodeengine/application/executors"
	nodeEngineApplicationMCP "github.com/blocknextai/platform-api/internal/nodeengine/application/mcp"
	"github.com/gofiber/fiber/v3"
)

type Dependencies struct {
	ServerURLTemplate string

	ServerService                         nodeEngineApplicationMCP.ServerService
	ExecutorService                       nodeEngineApplicationExecutors.ExecutorService
	CredentialOAuthTokenRegenerateService credentialOAuthApplicationRegenerate.CredentialOAuthTokenRegenerateService
}

type Module struct {
	handlers *mcpInfrastructure.Handlers
}

func NewModule(deps Dependencies) (*Module, error) {
	resolver := credentialresolver.NewCredentialResolver(deps.CredentialOAuthTokenRegenerateService)
	mcpAdapter := adapter.NewAdapter(deps.ExecutorService, resolver)

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
