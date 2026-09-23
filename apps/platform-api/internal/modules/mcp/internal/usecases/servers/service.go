package servers

import (
	"github.com/blocknextai/platform-api/internal/modules/mcp/internal/application/servers"
	mcpoauthContract "github.com/blocknextai/platform-api/internal/modules/mcpoauth/contract"
)

type Service struct {
	serverURLTemplate       string
	allServers              []*servers.Server
	mcpOAuthMetadataService mcpoauthContract.MetadataService
}

func NewService(
	serverURLTemplate string,
	allServers []*servers.Server,
	mcpOAuthMetadataService mcpoauthContract.MetadataService,
) *Service {
	return &Service{
		serverURLTemplate:       serverURLTemplate,
		allServers:              allServers,
		mcpOAuthMetadataService: mcpOAuthMetadataService,
	}
}
