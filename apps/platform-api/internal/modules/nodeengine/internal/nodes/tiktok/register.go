package tiktok

import (
	"github.com/blocknextai/platform-api/internal/filegateway"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/application/jsonschema"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/executors"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/functioncalling"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/mcp"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/nodes"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/nodes/tiktok/publishpost"
)

func Register(fileGateway filegateway.FileGateway) {
	nodeID := "tiktok"

	publishPostNodeID := nodeID + "_publish_post"
	publishPostNode := publishpost.NewTiktokPublishPostNode(publishPostNodeID)
	publishPostValidator := jsonschema.New[publishpost.TiktokPublishPostExecutorInput](publishPostNode.GetInputSchema())
	publishPostExecutor := publishpost.NewTiktokPublishPostExecutor(publishPostNodeID, publishPostValidator, fileGateway)

	nodes.RegisterNode(publishPostNode)
	executors.RegisterExecutor(publishPostExecutor)
	functioncalling.RegisterFunctionCalling(functioncalling.Generate(publishPostNode))

	mcp.RegisterServer(&mcp.Server{
		ID:          nodeID,
		Name:        "TikTok",
		Description: "Tools for publishing posts to TikTok.",
		Icon: mcp.ServerIcon{
			Brand: "tiktok",
		},
		Version: "0.0.1",
		Tools: []nodes.NodeManager{
			publishPostNode,
		},
	})
}
