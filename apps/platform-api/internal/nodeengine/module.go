package nodeengine

import (
	cachemiddleware "github.com/blocknextai/go-packages/fiber/middleware/cache"
	"github.com/blocknextai/platform-api/internal/filegateway"
	nodeEngineApplicationAdapters "github.com/blocknextai/platform-api/internal/nodeengine/application/adapters"
	nodeEngineApplicationCredentials "github.com/blocknextai/platform-api/internal/nodeengine/application/credentials"
	nodeEngineApplicationExecutors "github.com/blocknextai/platform-api/internal/nodeengine/application/executors"
	nodeEngineApplicationMCP "github.com/blocknextai/platform-api/internal/nodeengine/application/mcp"
	nodeEngineApplicationNodes "github.com/blocknextai/platform-api/internal/nodeengine/application/nodes"
	nodeEngineInfrastructure "github.com/blocknextai/platform-api/internal/nodeengine/infrastructure"
	nodeEnginePresentation "github.com/blocknextai/platform-api/internal/nodeengine/presentation"
	"github.com/gofiber/fiber/v3"
)

type Dependencies struct {
	FileGateway filegateway.FileGateway

	OAuth2RedirectURL         string
	WebhookTriggerURLTemplate string
}

type Module struct {
	NodeService         nodeEngineApplicationNodes.NodeService
	CredentialService   nodeEngineApplicationCredentials.CredentialService
	AdapterService      nodeEngineApplicationAdapters.AdapterService
	CredentialProcessor nodeEngineApplicationCredentials.CredentialProcessor
	ExecutorService     nodeEngineApplicationExecutors.ExecutorService
	MCPServerService    nodeEngineApplicationMCP.ServerService

	handlers *nodeEngineInfrastructure.Handlers
}

func NewModule(deps Dependencies) *Module {
	credentialService := nodeEngineApplicationCredentials.NewCredentialService()
	nodeService := nodeEngineApplicationNodes.NewNodeService()
	adapterService := nodeEngineApplicationAdapters.NewAdapterService()
	credentialProcessor := nodeEngineApplicationCredentials.NewCredentialProcessor(credentialService)
	executorService := nodeEngineApplicationExecutors.NewExecutorService()
	mcpServerService := nodeEngineApplicationMCP.NewServerService()

	handlers := nodeEngineInfrastructure.RegisterInfrastructure(nodeEngineInfrastructure.RegisterInfrastructureDeps{
		OAuth2RedirectURL:  deps.OAuth2RedirectURL,
		WebhookURLTemplate: deps.WebhookTriggerURLTemplate,

		NodeService:        nodeService,
		CredentialService:  credentialService,
		AdapterService:     adapterService,
		FileGatewayService: deps.FileGateway,
	})

	return &Module{
		NodeService:         nodeService,
		CredentialService:   credentialService,
		AdapterService:      adapterService,
		CredentialProcessor: credentialProcessor,
		ExecutorService:     executorService,
		MCPServerService:    mcpServerService,
		handlers:            handlers,
	}
}

func (m *Module) Register(router fiber.Router, cacheMiddleware *cachemiddleware.Middleware) {
	nodeEnginePresentation.RegisterPresentation(router, cacheMiddleware, m.handlers)
}
