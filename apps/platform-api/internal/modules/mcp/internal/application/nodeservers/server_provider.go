package nodeservers

import (
	"github.com/blocknextai/platform-api/internal/modules/mcp/internal/application/adapter"
	"github.com/blocknextai/platform-api/internal/modules/mcp/internal/application/servers"
	nodeengineContract "github.com/blocknextai/platform-api/internal/modules/nodeengine/contract"
)

const (
	MCPInvokeScope = "mcp:invoke"
)

type serverProvider struct {
	serverService nodeengineContract.ServerService
	adapter       adapter.Adapter
}

func NewServerProvider(
	serverService nodeengineContract.ServerService,
	adapter adapter.Adapter,
) servers.ServerProvider {
	return &serverProvider{
		serverService: serverService,
		adapter:       adapter,
	}
}

func (p *serverProvider) GetAllServers() ([]*servers.Server, error) {
	nodeServers := p.serverService.GetAllServers()
	result := make([]*servers.Server, 0, len(nodeServers))

	for _, nodeServer := range nodeServers {
		mcpServer, err := p.adapter.Build(nodeServer)
		if err != nil {
			return nil, err
		}

		tools := make([]servers.Tool, 0, len(nodeServer.GetTools()))
		for _, node := range nodeServer.GetTools() {
			icon := node.GetIcon()
			tools = append(tools, servers.Tool{
				ID:          node.GetID(),
				Version:     node.GetVersion(),
				Name:        node.GetName(),
				Description: node.GetDescription(),
				Icon: servers.Icon{
					Brand: icon.Brand,
					Glyph: icon.Glyph,
				},
				InputSchema:          node.GetInputSchema(),
				OutputSchema:         node.GetOutputSchema(),
				SupportedCredentials: node.GetSupportedCredentials(),
			})
		}

		result = append(result, &servers.Server{
			ID:          nodeServer.GetID(),
			Name:        nodeServer.GetName(),
			Description: nodeServer.GetDescription(),
			Icon: servers.Icon{
				Brand: nodeServer.GetIcon().Brand,
			},
			Version:      nodeServer.GetVersion(),
			Instructions: nodeServer.GetInstructions(),
			AuthMethods: []servers.AuthMethod{
				servers.AuthMethodAPIKey,
				servers.AuthMethodOAuth,
			},
			Scopes: []string{
				MCPInvokeScope,
			},
			Tools:     tools,
			MCPServer: mcpServer,
		})
	}

	return result, nil
}
