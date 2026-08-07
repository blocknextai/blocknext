package getallservers

import (
	nodeEngineDomainMCP "github.com/blocknextai/platform-api/internal/nodeengine/domain/mcp"
	nodeEngineDomainNodes "github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
)

type ServerResponse struct {
	ID           string                              `json:"id,omitempty"`
	Name         string                              `json:"name,omitempty"`
	Description  string                              `json:"description,omitempty"`
	Icon         nodeEngineDomainMCP.ServerIcon      `json:"icon,omitzero"`
	Version      string                              `json:"version,omitempty"`
	Instructions string                              `json:"instructions,omitempty"`
	URL          string                              `json:"url,omitempty"`
	Tools        []nodeEngineDomainNodes.NodeManager `json:"tools,omitempty"`
}

type GetAllServersResponse = []ServerResponse
