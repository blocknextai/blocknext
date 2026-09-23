package x

import (
	"github.com/blocknextai/platform-api/internal/filegateway"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/application/jsonschema"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/executors"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/functioncalling"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/mcp"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/nodes"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/nodes/x/publishmediapost"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/nodes/x/publishpost"
)

func Register(fileGateway filegateway.FileGateway) {
	nodeID := "x"

	publishPostNodeID := nodeID + "_publish_post"
	publishPostNode := publishpost.NewXPublishPostNode(publishPostNodeID)
	publishPostValidator := jsonschema.New[publishpost.XPublishPostExecutorInput](publishPostNode.GetInputSchema())
	publishPostExecutor := publishpost.NewXPublishPostExecutor(publishPostNodeID, publishPostValidator)

	nodes.RegisterNode(publishPostNode)
	executors.RegisterExecutor(publishPostExecutor)
	functioncalling.RegisterFunctionCalling(functioncalling.Generate(publishPostNode))

	publishMediaPostNodeID := nodeID + "_publish_media_post"
	publishMediaPostNode := publishmediapost.NewXPublishMediaPostNode(publishMediaPostNodeID)
	publishMediaPostValidator := jsonschema.New[publishmediapost.XPublishMediaPostExecutorInput](publishMediaPostNode.GetInputSchema())
	publishMediaPostExecutor := publishmediapost.NewXPublishMediaPostExecutor(publishMediaPostNodeID, publishMediaPostValidator, fileGateway)

	nodes.RegisterNode(publishMediaPostNode)
	executors.RegisterExecutor(publishMediaPostExecutor)
	functioncalling.RegisterFunctionCalling(functioncalling.Generate(publishMediaPostNode))

	mcp.RegisterServer(&mcp.Server{
		ID:          nodeID,
		Name:        "X",
		Description: "Tools for publishing posts and media to X (Twitter).",
		Icon: mcp.ServerIcon{
			Brand: "x",
		},
		Version: "0.0.1",
		Tools: []nodes.NodeManager{
			publishPostNode,
			publishMediaPostNode,
		},
	})
}
