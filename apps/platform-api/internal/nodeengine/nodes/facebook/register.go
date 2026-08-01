package facebook

import (
	"github.com/blocknextai/platform-api/internal/nodeengine/application/jsonschema"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/executors"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/functioncalling"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/mcp"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/facebook/publishmediapost"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/facebook/publishpost"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/facebook/publishstory"
)

func Register() {
	nodeID := "facebook"

	publishPostNodeID := nodeID + "_publish_post"
	publishPostNode := publishpost.NewFacebookPublishPostNode(publishPostNodeID)
	publishPostValidator := jsonschema.New[publishpost.FacebookPublishPostExecutorInput](publishPostNode.GetInputSchema())
	publishPostExecutor := publishpost.NewFacebookPublishPostExecutor(publishPostNodeID, publishPostValidator)

	nodes.RegisterNode(publishPostNode)
	executors.RegisterExecutor(publishPostExecutor)
	functioncalling.RegisterFunctionCalling(functioncalling.Generate(publishPostNode))

	publishMediaPostNodeID := nodeID + "_publish_media_post"
	publishMediaPostNode := publishmediapost.NewFacebookPublishMediaPostNode(publishMediaPostNodeID)
	publishMediaPostValidator := jsonschema.New[publishmediapost.FacebookPublishMediaPostExecutorInput](publishMediaPostNode.GetInputSchema())
	publishMediaPostExecutor := publishmediapost.NewFacebookPublishMediaPostExecutor(publishMediaPostNodeID, publishMediaPostValidator)

	nodes.RegisterNode(publishMediaPostNode)
	executors.RegisterExecutor(publishMediaPostExecutor)
	functioncalling.RegisterFunctionCalling(functioncalling.Generate(publishMediaPostNode))

	publishStoryNodeID := nodeID + "_publish_story"
	publishStoryNode := publishstory.NewFacebookPublishStoryNode(publishStoryNodeID)
	publishStoryValidator := jsonschema.New[publishstory.FacebookPublishStoryExecutorInput](publishStoryNode.GetInputSchema())
	publishStoryExecutor := publishstory.NewFacebookPublishStoryExecutor(publishStoryNodeID, publishStoryValidator)

	nodes.RegisterNode(publishStoryNode)
	executors.RegisterExecutor(publishStoryExecutor)
	functioncalling.RegisterFunctionCalling(functioncalling.Generate(publishStoryNode))

	mcp.RegisterServer(&mcp.Server{
		ID:          nodeID,
		Name:        "Facebook",
		Description: "Tools for publishing posts, media, and stories to Facebook pages.",
		Version:     "0.0.1",
		Tools: []nodes.NodeManager{
			publishPostNode,
			publishMediaPostNode,
			publishStoryNode,
		},
	})
}
