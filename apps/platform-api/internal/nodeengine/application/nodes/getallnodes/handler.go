package getallnodes

import (
	"context"

	nodeEngineApplicationNodes "github.com/blocknextai/platform-api/internal/nodeengine/application/nodes"
)

type Handler struct {
	nodeService nodeEngineApplicationNodes.NodeService
}

func New(
	nodeService nodeEngineApplicationNodes.NodeService,
) *Handler {
	return &Handler{
		nodeService: nodeService,
	}
}

func (h *Handler) Handle(_ context.Context, _ *GetAllNodesQuery) (*GetAllNodesResponse, error) {
	return new(h.nodeService.GetAllNodes()), nil
}
