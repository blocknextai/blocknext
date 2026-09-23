package contract

import (
	nodeengineApplicationNodes "github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/application/nodes"
	nodeengineDomainNodes "github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/nodes"
)

type NodeAnnotations = nodeengineDomainNodes.NodeAnnotations
type NodeManager = nodeengineDomainNodes.NodeManager

type NodeService = nodeengineApplicationNodes.NodeService

var GetNode = nodeengineDomainNodes.GetNode
