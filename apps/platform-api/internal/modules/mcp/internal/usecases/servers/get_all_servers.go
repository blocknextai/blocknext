package servers

import (
	"context"
	"strings"

	gjs "github.com/google/jsonschema-go/jsonschema"

	"github.com/blocknextai/platform-api/internal/modules/mcp/internal/application/servers"
)

type GetAllServersQuery struct{}

type IconResponse struct {
	Brand string `json:"brand,omitempty"`
	Glyph string `json:"glyph,omitempty"`
}

type ToolResponse struct {
	ID                   string       `json:"id,omitempty"`
	Version              string       `json:"version,omitempty"`
	Name                 string       `json:"name,omitempty"`
	Description          string       `json:"description,omitempty"`
	Icon                 IconResponse `json:"icon,omitzero"`
	InputSchema          *gjs.Schema  `json:"inputSchema,omitempty"`
	OutputSchema         *gjs.Schema  `json:"outputSchema,omitempty"`
	SupportedCredentials []string     `json:"supportedCredentials,omitempty"`
	Scopes               []string     `json:"scopes,omitempty"`
}

type ServerResponse struct {
	ID           string               `json:"id,omitempty"`
	AuthMethods  []servers.AuthMethod `json:"authMethods,omitempty"`
	Name         string               `json:"name,omitempty"`
	Description  string               `json:"description,omitempty"`
	Icon         IconResponse         `json:"icon,omitzero"`
	Version      string               `json:"version,omitempty"`
	Instructions string               `json:"instructions,omitempty"`
	URL          string               `json:"url,omitempty"`
	Tools        []ToolResponse       `json:"tools,omitempty"`
}

type GetAllServersResponse = []ServerResponse

func (s *Service) GetAllServers(ctx context.Context, _ *GetAllServersQuery) (*GetAllServersResponse, error) {
	return MapServersToGetAllServersResponse(s.allServers, s.serverURLTemplate), nil
}

func MapServersToGetAllServersResponse(
	allServers []*servers.Server,
	serverURLTemplate string,
) *GetAllServersResponse {
	response := make(GetAllServersResponse, 0, len(allServers))
	for _, server := range allServers {
		response = append(response, ServerResponse{
			ID:           server.ID,
			AuthMethods:  server.AuthMethods,
			Name:         server.Name,
			Description:  server.Description,
			Icon:         MapIconToIconResponse(server.Icon),
			Version:      server.Version,
			Instructions: server.Instructions,
			URL:          strings.ReplaceAll(serverURLTemplate, "{serverId}", server.ID),
			Tools:        MapToolsToToolResponses(server.Tools),
		})
	}
	return &response
}

func MapToolsToToolResponses(tools []servers.Tool) []ToolResponse {
	responses := make([]ToolResponse, 0, len(tools))
	for _, tool := range tools {
		responses = append(responses, ToolResponse{
			ID:                   tool.ID,
			Version:              tool.Version,
			Name:                 tool.Name,
			Description:          tool.Description,
			Icon:                 MapIconToIconResponse(tool.Icon),
			InputSchema:          tool.InputSchema,
			OutputSchema:         tool.OutputSchema,
			SupportedCredentials: tool.SupportedCredentials,
			Scopes:               tool.Scopes,
		})
	}
	return responses
}

func MapIconToIconResponse(icon servers.Icon) IconResponse {
	return IconResponse{
		Brand: icon.Brand,
		Glyph: icon.Glyph,
	}
}
