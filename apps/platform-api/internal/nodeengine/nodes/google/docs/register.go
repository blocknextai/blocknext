package docs

import (
	"github.com/blocknextai/platform-api/internal/nodeengine/application/jsonschema"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/executors"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/functioncalling"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/mcp"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/google/docs/createorupdate"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/google/docs/readdocs"
)

func Register() {
	nodeID := "google_docs"

	createOrUpdateNodeID := nodeID + "_create_or_update"
	createOrUpdateNode := createorupdate.NewGoogleDocsCreateOrUpdateNode(createOrUpdateNodeID)
	createOrUpdateValidator := jsonschema.New[createorupdate.GoogleDocsCreateOrUpdateExecutorInput](createOrUpdateNode.GetInputSchema())
	createOrUpdateExecutor := createorupdate.NewGoogleDocsCreateOrUpdateExecutor(createOrUpdateNodeID, createOrUpdateValidator)

	nodes.RegisterNode(createOrUpdateNode)
	executors.RegisterExecutor(createOrUpdateExecutor)
	functioncalling.RegisterFunctionCalling(functioncalling.Generate(createOrUpdateNode))

	readDocsNodeID := nodeID + "_read_docs"
	readDocsNode := readdocs.NewGoogleDocsReadDocsNode(readDocsNodeID)
	readDocsValidator := jsonschema.New[readdocs.GoogleDocsReadDocsExecutorInput](readDocsNode.GetInputSchema())
	readDocsExecutor := readdocs.NewGoogleDocsReadDocsExecutor(readDocsNodeID, readDocsValidator)

	nodes.RegisterNode(readDocsNode)
	executors.RegisterExecutor(readDocsExecutor)
	functioncalling.RegisterFunctionCalling(functioncalling.Generate(readDocsNode))

	mcp.RegisterServer(&mcp.Server{
		ID:          nodeID,
		Name:        "Google Docs",
		Description: "Tools for creating, updating, and reading Google Docs.",
		Icon: mcp.ServerIcon{
			Brand: "google_docs",
		},
		Version: "0.0.1",
		Tools: []nodes.NodeManager{
			createOrUpdateNode,
			readDocsNode,
		},
	})
}
