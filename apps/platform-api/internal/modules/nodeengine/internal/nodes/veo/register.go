package veo

import (
	"github.com/blocknextai/platform-api/internal/filegateway"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/application/jsonschema"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/mcp"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/providers"
)

func Register(fileGatewayService filegateway.FileGateway) {
	nodeID := "veo"

	veoNode := NewVeoNode(nodeID)
	veoValidator := jsonschema.New[VeoExecutorInput](veoNode.GetInputSchema())
	veoExecutor := NewVeoExecutor(nodeID, veoValidator, fileGatewayService)

	providers.Register(providers.Provider{
		Server: mcp.Server{
			ID:          nodeID,
			Name:        "Veo",
			Description: "Tools for video generation via Google Veo.",
			Icon: mcp.ServerIcon{
				Brand: "veo",
			},
			Version: "0.0.1",
		},
		Actions: []providers.Action{
			{Node: veoNode, Executor: veoExecutor},
		},
	})
}
