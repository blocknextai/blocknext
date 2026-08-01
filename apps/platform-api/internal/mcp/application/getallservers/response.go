package getallservers

import (
	nodeEngineDomainNodes "github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
)

type ServerResponse struct {
	ID           string                              `json:"id,omitempty"`
	Name         string                              `json:"name,omitempty"`
	Description  string                              `json:"description,omitempty"`
	Version      string                              `json:"version,omitempty"`
	Instructions string                              `json:"instructions,omitempty"`
	URL          string                              `json:"url,omitempty"`
	Tools        []nodeEngineDomainNodes.NodeManager `json:"tools,omitempty"`
}

type GetAllServersResponse = []ServerResponse
