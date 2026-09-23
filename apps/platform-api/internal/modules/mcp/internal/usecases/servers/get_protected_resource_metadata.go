package servers

import (
	"context"
	"slices"

	"github.com/blocknextai/platform-api/internal/modules/mcp/internal/application/servers"
	mcpoauthContract "github.com/blocknextai/platform-api/internal/modules/mcpoauth/contract"
)

type GetProtectedResourceMetadataQuery struct {
	ServerID     string
	ResourcePath string
}

type GetProtectedResourceMetadataResponse struct {
	Resource               string   `json:"resource"`
	AuthorizationServers   []string `json:"authorization_servers"`
	ScopesSupported        []string `json:"scopes_supported"`
	BearerMethodsSupported []string `json:"bearer_methods_supported"`
	ResourceName           string   `json:"resource_name"`
}

func (s *Service) GetProtectedResourceMetadata(ctx context.Context, query *GetProtectedResourceMetadataQuery) (*GetProtectedResourceMetadataResponse, error) {
	index := slices.IndexFunc(s.allServers, func(server *servers.Server) bool {
		return server.ID == query.ServerID
	})
	if index < 0 {
		return nil, servers.ErrServerNotFound
	}

	server := s.allServers[index]
	metadata := s.mcpOAuthMetadataService.GetProtectedResourceMetadata(query.ResourcePath, server.Scopes)

	return MapProtectedResourceMetadataToResponse(metadata, server.Name), nil
}

func MapProtectedResourceMetadataToResponse(
	metadata *mcpoauthContract.ProtectedResourceMetadata,
	resourceName string,
) *GetProtectedResourceMetadataResponse {
	return &GetProtectedResourceMetadataResponse{
		Resource:               metadata.Resource,
		AuthorizationServers:   metadata.AuthorizationServers,
		ScopesSupported:        metadata.ScopesSupported,
		BearerMethodsSupported: metadata.BearerMethodsSupported,
		ResourceName:           resourceName,
	}
}
