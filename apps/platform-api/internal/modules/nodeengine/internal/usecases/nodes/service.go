package nodes

import (
	nodeEngineApplicationNodes "github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/application/nodes"
)

type Service struct {
	nodeService nodeEngineApplicationNodes.NodeService
}

func NewService(
	nodeService nodeEngineApplicationNodes.NodeService,
) *Service {
	return &Service{
		nodeService: nodeService,
	}
}
