package usecases

import (
	"errors"
	"net/http"
	"slices"

	serversUseCases "github.com/blocknextai/platform-api/internal/modules/mcp/internal/usecases/servers"

	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/blocknextai/platform-api/internal/modules/mcp/internal/application/servers"
	mcpoauthContract "github.com/blocknextai/platform-api/internal/modules/mcpoauth/contract"
)

type ServerHandler struct {
	ID          string
	AuthMethods []servers.AuthMethod
	Scopes      []string
	Handler     http.Handler
}

type Services struct {
	Servers        *serversUseCases.Service
	ServerHandlers []*ServerHandler
}

type ServiceDependencies struct {
	ServerURLTemplate string

	ServerProviders         []servers.ServerProvider
	MCPOAuthMetadataService mcpoauthContract.MetadataService
}

func NewServices(deps ServiceDependencies) (*Services, error) {
	var allServers []*servers.Server
	for _, serverProvider := range deps.ServerProviders {
		providedServers, err := serverProvider.GetAllServers()
		if err != nil {
			return nil, err
		}

		allServers = append(allServers, providedServers...)
	}

	serverHandlers := make([]*ServerHandler, 0, len(allServers))
	for _, server := range allServers {
		if slices.ContainsFunc(serverHandlers, func(serverHandler *ServerHandler) bool {
			return serverHandler.ID == server.ID
		}) {
			return nil, servers.ErrDuplicateServerID.WithCause(errors.New(server.ID))
		}

		serverHandlers = append(serverHandlers, &ServerHandler{
			ID:          server.ID,
			AuthMethods: server.AuthMethods,
			Scopes:      server.Scopes,
			Handler: mcpsdk.NewStreamableHTTPHandler(
				func(*http.Request) *mcpsdk.Server { return server.MCPServer },
				&mcpsdk.StreamableHTTPOptions{
					Stateless: true,
				},
			),
		})
	}

	return &Services{
		Servers:        serversUseCases.NewService(deps.ServerURLTemplate, allServers, deps.MCPOAuthMetadataService),
		ServerHandlers: serverHandlers,
	}, nil
}
