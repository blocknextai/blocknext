package mcp

import (
	"sync"

	"github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
)

var (
	once     sync.Once
	instance *ServerRegistry
)

type ServerRegistry struct {
	serversMap map[string]ServerManager
	servers    []ServerManager
}

func GetRegistry() *ServerRegistry {
	once.Do(func() {
		instance = &ServerRegistry{
			serversMap: make(map[string]ServerManager),
			servers:    make([]ServerManager, 0),
		}
	})
	return instance
}

func (r *ServerRegistry) RegisterServer(server ServerManager) {
	if server.GetDisabled() {
		return
	}

	server.SetTools(filterEnabledTools(server.GetTools()))

	r.serversMap[server.GetID()] = server
	r.servers = append(r.servers, server)
}

func filterEnabledTools(tools []nodes.NodeManager) []nodes.NodeManager {
	enabled := make([]nodes.NodeManager, 0, len(tools))
	for _, tool := range tools {
		if tool.GetDisabled() {
			continue
		}
		enabled = append(enabled, tool)
	}
	return enabled
}

func (r *ServerRegistry) GetServer(id string) (ServerManager, bool) {
	server, ok := r.serversMap[id]
	return server, ok
}

func (r *ServerRegistry) GetAllServers() []ServerManager {
	return r.servers
}

func RegisterServer(server ServerManager) {
	GetRegistry().RegisterServer(server)
}

func GetServer(id string) (ServerManager, bool) {
	return GetRegistry().GetServer(id)
}

func GetAllServers() []ServerManager {
	return GetRegistry().GetAllServers()
}
