package getallservers

import (
	"strings"

	nodeEngineDomainMCP "github.com/blocknextai/platform-api/internal/nodeengine/domain/mcp"
)

func MapGetAllServersQueryToGetAllServersResponse(
	servers []nodeEngineDomainMCP.ServerManager,
	serverURLTemplate string,
) *GetAllServersResponse {
	response := make(GetAllServersResponse, 0, len(servers))
	for _, server := range servers {
		url := strings.ReplaceAll(serverURLTemplate, "{serverId}", server.GetID())
		response = append(response, ServerResponse{
			ID:           server.GetID(),
			Name:         server.GetName(),
			Description:  server.GetDescription(),
			Icon:         server.GetIcon(),
			Version:      server.GetVersion(),
			Instructions: server.GetInstructions(),
			URL:          url,
			Tools:        server.GetTools(),
		})
	}
	return &response
}
