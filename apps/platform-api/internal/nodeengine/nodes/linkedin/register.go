package linkedin

import (
	"github.com/blocknextai/platform-api/internal/filegateway"
	"github.com/blocknextai/platform-api/internal/nodeengine/application/jsonschema"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/executors"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/functioncalling"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/mcp"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/linkedin/publishmediapost"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/linkedin/publishpost"
)

func Register(fileGateway filegateway.FileGateway) {
	nodeID := "linkedin"

	publishPostNodeID := nodeID + "_publish_post"
	publishPostNode := publishpost.NewLinkedinPublishPostNode(publishPostNodeID)
	publishPostValidator := jsonschema.New[publishpost.LinkedinPublishPostExecutorInput](publishPostNode.GetInputSchema())
	publishPostExecutor := publishpost.NewLinkedinPublishPostExecutor(publishPostNodeID, publishPostValidator)

	nodes.RegisterNode(publishPostNode)
	executors.RegisterExecutor(publishPostExecutor)
	functioncalling.RegisterFunctionCalling(functioncalling.Generate(publishPostNode))

	publishMediaPostNodeID := nodeID + "_publish_media_post"
	publishMediaPostNode := publishmediapost.NewLinkedinPublishMediaPostNode(publishMediaPostNodeID)
	publishMediaPostValidator := jsonschema.New[publishmediapost.LinkedinPublishMediaPostExecutorInput](publishMediaPostNode.GetInputSchema())
	publishMediaPostExecutor := publishmediapost.NewLinkedinPublishMediaPostExecutor(publishMediaPostNodeID, publishMediaPostValidator, fileGateway)

	nodes.RegisterNode(publishMediaPostNode)
	executors.RegisterExecutor(publishMediaPostExecutor)
	functioncalling.RegisterFunctionCalling(functioncalling.Generate(publishMediaPostNode))

	mcp.RegisterServer(&mcp.Server{
		ID:          nodeID,
		Name:        "LinkedIn",
		Description: "Tools for publishing posts and media to LinkedIn.",
		Icon: mcp.ServerIcon{
			Brand: "linkedin",
		},
		Version: "0.0.1",
		Tools: []nodes.NodeManager{
			publishPostNode,
			publishMediaPostNode,
		},
	})
}
