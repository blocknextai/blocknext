package usecases

import (
	"github.com/blocknextai/platform-api/internal/filegateway"
	nodeEngineApplicationAdapters "github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/application/adapters"
	nodeEngineApplicationCredentials "github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/application/credentials"
	nodeEngineApplicationNodes "github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/application/nodes"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/credentials"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/nodes"
	adaptersUseCases "github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/usecases/adapters"
	credentialsUseCases "github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/usecases/credentials"
	nodesUseCases "github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/usecases/nodes"
)

type Services struct {
	Adapters    *adaptersUseCases.Service
	Credentials *credentialsUseCases.Service
	Nodes       *nodesUseCases.Service
}

type ServiceDependencies struct {
	OAuth2RedirectURL  string
	WebhookURLTemplate string

	NodeService        nodeEngineApplicationNodes.NodeService
	CredentialService  nodeEngineApplicationCredentials.CredentialService
	AdapterService     nodeEngineApplicationAdapters.AdapterService
	FileGatewayService filegateway.FileGateway
}

func NewServices(deps ServiceDependencies) *Services {
	credentials.Register(deps.OAuth2RedirectURL)
	nodes.Register(deps.FileGatewayService)

	return &Services{
		Adapters:    adaptersUseCases.NewService(deps.AdapterService, deps.WebhookURLTemplate),
		Credentials: credentialsUseCases.NewService(deps.CredentialService),
		Nodes:       nodesUseCases.NewService(deps.NodeService),
	}
}
