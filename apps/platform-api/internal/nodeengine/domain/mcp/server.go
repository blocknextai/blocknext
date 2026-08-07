package mcp

import (
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
)

type Server struct {
	ID           string              `json:"id,omitempty"`
	Name         string              `json:"name,omitempty"`
	Description  string              `json:"description,omitempty"`
	Icon         ServerIcon          `json:"icon,omitzero"`
	Version      string              `json:"version,omitempty"`
	Instructions string              `json:"instructions,omitempty"`
	Tools        []nodes.NodeManager `json:"tools,omitempty"`
	Disabled     bool                `json:"disabled,omitempty"`
}

type ServerManager interface {
	GetID() string
	GetName() string
	GetDescription() string
	GetIcon() ServerIcon
	GetVersion() string
	GetInstructions() string
	GetTools() []nodes.NodeManager
	SetTools(tools []nodes.NodeManager)
	GetDisabled() bool
}

func (s *Server) GetID() string {
	return s.ID
}

func (s *Server) GetName() string {
	return s.Name
}

func (s *Server) GetDescription() string {
	return s.Description
}

func (s *Server) GetIcon() ServerIcon {
	return s.Icon
}

func (s *Server) GetVersion() string {
	return s.Version
}

func (s *Server) GetInstructions() string {
	return s.Instructions
}

func (s *Server) GetTools() []nodes.NodeManager {
	return s.Tools
}

func (s *Server) SetTools(tools []nodes.NodeManager) {
	s.Tools = tools
}

func (s *Server) GetDisabled() bool {
	return s.Disabled
}
