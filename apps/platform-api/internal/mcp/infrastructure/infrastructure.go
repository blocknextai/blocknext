package infrastructure

import (
	"net/http"

	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
	"github.com/blocknextai/platform-api/internal/mcp/application/adapter"
	"github.com/blocknextai/platform-api/internal/mcp/application/getallservers"
	nodeEngineApplicationMCP "github.com/blocknextai/platform-api/internal/nodeengine/application/mcp"
	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

type ServerHandler struct {
	ID      string
	Name    string
	Handler http.Handler
}

type Handlers struct {
	GetAllServers  cqrs.Handler[*getallservers.GetAllServersQuery, *getallservers.GetAllServersResponse]
	ServerHandlers []*ServerHandler
}

type RegisterInfrastructureDeps struct {
	ServerURLTemplate string

	ServerService nodeEngineApplicationMCP.ServerService
	Adapter       adapter.Adapter
}

func RegisterInfrastructure(deps RegisterInfrastructureDeps) (*Handlers, error) {
	servers := deps.ServerService.GetAllServers()
	serverHandlers := make([]*ServerHandler, 0, len(servers))

	for _, server := range servers {
		mcpServer, err := deps.Adapter.Build(server)
		if err != nil {
			return nil, err
		}

		boundServer := mcpServer
		streamableHandler := mcpsdk.NewStreamableHTTPHandler(
			func(*http.Request) *mcpsdk.Server { return boundServer },
			&mcpsdk.StreamableHTTPOptions{
				Stateless: true,
			},
		)

		serverHandlers = append(serverHandlers, &ServerHandler{
			ID:      server.GetID(),
			Name:    server.GetName(),
			Handler: streamableHandler,
		})
	}

	return &Handlers{
		GetAllServers:  cqrs.ValidationBehavior(getallservers.New(deps.ServerURLTemplate, deps.ServerService)),
		ServerHandlers: serverHandlers,
	}, nil
}
