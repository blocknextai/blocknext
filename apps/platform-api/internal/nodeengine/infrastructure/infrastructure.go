package infrastructure

import (
	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
	"github.com/blocknextai/platform-api/internal/filegateway"
	nodeEngineApplicationAdapters "github.com/blocknextai/platform-api/internal/nodeengine/application/adapters"
	"github.com/blocknextai/platform-api/internal/nodeengine/application/adapters/getalltriggervariables"
	"github.com/blocknextai/platform-api/internal/nodeengine/application/adapters/getallwebhooksources"
	nodeEngineApplicationCredentials "github.com/blocknextai/platform-api/internal/nodeengine/application/credentials"
	"github.com/blocknextai/platform-api/internal/nodeengine/application/credentials/getallcredentials"
	"github.com/blocknextai/platform-api/internal/nodeengine/application/credentials/getcredentialbyid"
	nodeEngineApplicationNodes "github.com/blocknextai/platform-api/internal/nodeengine/application/nodes"
	"github.com/blocknextai/platform-api/internal/nodeengine/application/nodes/getallnodes"
	"github.com/blocknextai/platform-api/internal/nodeengine/credentials"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes"
)

type Handlers struct {
	GetAllNodes            cqrs.Handler[*getallnodes.GetAllNodesQuery, *getallnodes.GetAllNodesResponse]
	GetAllCredentials      cqrs.Handler[*getallcredentials.GetAllCredentialsQuery, *getallcredentials.GetAllCredentialsResponse]
	GetCredentialByID      cqrs.Handler[*getcredentialbyid.GetCredentialByIDQuery, *getcredentialbyid.GetCredentialByIDResponse]
	GetAllWebhookSources   cqrs.Handler[*getallwebhooksources.GetAllWebhookSourcesQuery, *getallwebhooksources.GetAllWebhookSourcesResponse]
	GetAllTriggerVariables cqrs.Handler[*getalltriggervariables.GetAllTriggerVariablesQuery, *getalltriggervariables.GetAllTriggerVariablesResponse]
}

type RegisterInfrastructureDeps struct {
	OAuth2RedirectURL  string
	WebhookURLTemplate string

	NodeService        nodeEngineApplicationNodes.NodeService
	CredentialService  nodeEngineApplicationCredentials.CredentialService
	AdapterService     nodeEngineApplicationAdapters.AdapterService
	FileGatewayService filegateway.FileGateway
}

func RegisterInfrastructure(deps RegisterInfrastructureDeps) *Handlers {
	credentials.Register(deps.OAuth2RedirectURL)
	nodes.Register(deps.FileGatewayService)

	return &Handlers{
		GetAllNodes:            cqrs.ValidationBehavior(getallnodes.New(deps.NodeService)),
		GetAllCredentials:      cqrs.ValidationBehavior(getallcredentials.New(deps.CredentialService)),
		GetCredentialByID:      cqrs.ValidationBehavior(getcredentialbyid.New(deps.CredentialService)),
		GetAllWebhookSources:   cqrs.ValidationBehavior(getallwebhooksources.New(deps.AdapterService, deps.WebhookURLTemplate)),
		GetAllTriggerVariables: cqrs.ValidationBehavior(getalltriggervariables.New(deps.AdapterService)),
	}
}
