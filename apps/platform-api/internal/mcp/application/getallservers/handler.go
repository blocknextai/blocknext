package getallservers

import (
	"context"

	nodeEngineApplicationMCP "github.com/blocknextai/platform-api/internal/nodeengine/application/mcp"
)

type Handler struct {
	serverURLTemplate string
	serverService     nodeEngineApplicationMCP.ServerService
}

func New(
	serverURLTemplate string,
	serverService nodeEngineApplicationMCP.ServerService,
) *Handler {
	return &Handler{
		serverService:     serverService,
		serverURLTemplate: serverURLTemplate,
	}
}

func (h *Handler) Handle(_ context.Context, _ *GetAllServersQuery) (*GetAllServersResponse, error) {
	servers := h.serverService.GetAllServers()
	return MapGetAllServersQueryToGetAllServersResponse(servers, h.serverURLTemplate), nil
}
