package servers

import (
	gjs "github.com/google/jsonschema-go/jsonschema"
	mcpsdk "github.com/modelcontextprotocol/go-sdk/mcp"
)

type Icon struct {
	Brand string
	Glyph string
}

type Tool struct {
	ID                   string
	Version              string
	Name                 string
	Description          string
	Icon                 Icon
	InputSchema          *gjs.Schema
	OutputSchema         *gjs.Schema
	SupportedCredentials []string
	Scopes               []string
}

type Server struct {
	ID           string
	Name         string
	Description  string
	Icon         Icon
	Version      string
	Instructions string
	AuthMethods  []AuthMethod
	Scopes       []string
	Tools        []Tool
	MCPServer    *mcpsdk.Server
}

type ServerProvider interface {
	GetAllServers() ([]*Server, error)
}
