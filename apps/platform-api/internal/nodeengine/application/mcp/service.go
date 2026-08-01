package mcp

import (
	nodeEngineDomainMCP "github.com/blocknextai/platform-api/internal/nodeengine/domain/mcp"
)

type ServerService interface {
	GetAllServers() []nodeEngineDomainMCP.ServerManager
	GetServerByID(id string) (nodeEngineDomainMCP.ServerManager, bool)
}

type serverService struct{}

func NewServerService() ServerService {
	return &serverService{}
}

func (s *serverService) GetAllServers() []nodeEngineDomainMCP.ServerManager {
	return nodeEngineDomainMCP.GetAllServers()
}

func (s *serverService) GetServerByID(id string) (nodeEngineDomainMCP.ServerManager, bool) {
	return nodeEngineDomainMCP.GetServer(id)
}
