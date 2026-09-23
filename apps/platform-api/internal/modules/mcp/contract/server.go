package contract

import (
	mcpApplicationServers "github.com/blocknextai/platform-api/internal/modules/mcp/internal/application/servers"
)

type AuthMethod = mcpApplicationServers.AuthMethod

const AuthMethodOAuth = mcpApplicationServers.AuthMethodOAuth

type Server = mcpApplicationServers.Server
type ServerProvider = mcpApplicationServers.ServerProvider
type Tool = mcpApplicationServers.Tool
