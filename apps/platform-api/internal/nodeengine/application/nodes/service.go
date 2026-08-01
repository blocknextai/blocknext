package nodes

import (
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
)

type NodeService interface {
	GetAllNodes() []nodes.NodeManager
	GetNodeByID(id string) (nodes.NodeManager, bool)
	GetNodeSchemasByIDs(nodeIDs []string) []nodes.NodeManager
}

type nodeService struct{}

func NewNodeService() NodeService {
	return &nodeService{}
}

func (s *nodeService) GetAllNodes() []nodes.NodeManager {
	return nodes.GetAllNodes()
}

func (s *nodeService) GetNodeByID(id string) (nodes.NodeManager, bool) {
	return nodes.GetNode(id)
}

func (s *nodeService) GetNodeSchemasByIDs(nodeIDs []string) []nodes.NodeManager {
	nodeSchemasByKey := make(map[string]nodes.NodeManager)
	nodeSchemas := make([]nodes.NodeManager, 0)

	for _, nodeID := range nodeIDs {
		if _, exists := nodeSchemasByKey[nodeID]; exists {
			continue
		}

		node, exists := nodes.GetNode(nodeID)
		if !exists {
			return nil
		}

		nodeSchemasByKey[nodeID] = node
		nodeSchemas = append(nodeSchemas, node)
	}

	return nodeSchemas
}
