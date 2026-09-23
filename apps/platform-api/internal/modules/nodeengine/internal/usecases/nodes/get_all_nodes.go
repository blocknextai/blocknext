package nodes

import (
	"context"

	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/nodes"
)

type GetAllNodesQuery struct{}

type GetAllNodesResponse = []nodes.NodeManager

func (s *Service) GetAllNodes(ctx context.Context, _ *GetAllNodesQuery) (*GetAllNodesResponse, error) {
	return new(s.nodeService.GetAllNodes()), nil
}
