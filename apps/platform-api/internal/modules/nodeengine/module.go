package nodeengine

import (
	"github.com/gofiber/fiber/v3"

	cachemiddleware "github.com/blocknextai/go-packages/fiber/middleware/cache"
	"github.com/blocknextai/platform-api/internal/filegateway"
	nodeEngineApplicationAdapters "github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/application/adapters"
	nodeEngineApplicationCredentials "github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/application/credentials"
	nodeEngineApplicationExecutors "github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/application/executors"
	nodeEngineApplicationMCP "github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/application/mcp"
	nodeEngineApplicationNodes "github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/application/nodes"
	nodeEngineHTTP "github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/http"
	nodeEngineUseCases "github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/usecases"
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

	useCases *nodeEngineUseCases.Services
}

func NewModule(deps Dependencies) *Module {
	credentialService := nodeEngineApplicationCredentials.NewCredentialService()
	nodeService := nodeEngineApplicationNodes.NewNodeService()
	adapterService := nodeEngineApplicationAdapters.NewAdapterService()
	credentialProcessor := nodeEngineApplicationCredentials.NewCredentialProcessor(credentialService)
	executorService := nodeEngineApplicationExecutors.NewExecutorService()
	mcpServerService := nodeEngineApplicationMCP.NewServerService()

	useCases := nodeEngineUseCases.NewServices(nodeEngineUseCases.ServiceDependencies{
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
		useCases:            useCases,
	}
}

func (m *Module) Register(router fiber.Router, cacheMiddleware *cachemiddleware.Middleware) {
	nodeEngineHTTP.RegisterRoutes(router, cacheMiddleware, m.useCases)
}
