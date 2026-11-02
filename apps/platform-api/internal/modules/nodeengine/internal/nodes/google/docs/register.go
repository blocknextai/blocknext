package docs

import (
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/application/jsonschema"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/mcp"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/providers"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/nodes/google/docs/createorupdate"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/nodes/google/docs/readdocs"
)

func Register() {
	nodeID := "google_docs"

	createOrUpdateNodeID := nodeID + "_create_or_update"
	createOrUpdateNode := createorupdate.NewGoogleDocsCreateOrUpdateNode(createOrUpdateNodeID)
	createOrUpdateValidator := jsonschema.New[createorupdate.GoogleDocsCreateOrUpdateExecutorInput](createOrUpdateNode.GetInputSchema())
	createOrUpdateExecutor := createorupdate.NewGoogleDocsCreateOrUpdateExecutor(createOrUpdateNodeID, createOrUpdateValidator)

	readDocsNodeID := nodeID + "_read_docs"
	readDocsNode := readdocs.NewGoogleDocsReadDocsNode(readDocsNodeID)
	readDocsValidator := jsonschema.New[readdocs.GoogleDocsReadDocsExecutorInput](readDocsNode.GetInputSchema())
	readDocsExecutor := readdocs.NewGoogleDocsReadDocsExecutor(readDocsNodeID, readDocsValidator)

	providers.Register(providers.Provider{
		Server: mcp.Server{
			ID:          nodeID,
			Name:        "Google Docs",
			Description: "Tools for creating, updating, and reading Google Docs.",
			Icon: mcp.ServerIcon{
				Brand: "google_docs",
			},
			Version: "0.0.1",
		},
		Actions: []providers.Action{
			{Node: createOrUpdateNode, Executor: createOrUpdateExecutor},
			{Node: readDocsNode, Executor: readDocsExecutor},
		},
	})
}
