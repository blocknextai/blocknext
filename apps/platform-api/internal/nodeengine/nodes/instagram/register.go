package instagram

import (
	"github.com/blocknextai/platform-api/internal/nodeengine/application/jsonschema"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/executors"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/functioncalling"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/mcp"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/instagram/publishpost"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/instagram/publishreels"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/instagram/publishstory"
)

func Register() {
	nodeID := "instagram"

	publishPostNodeID := nodeID + "_publish_post"
	publishPostNode := publishpost.NewInstagramPublishPostNode(publishPostNodeID)
	publishPostValidator := jsonschema.New[publishpost.InstagramPublishPostExecutorInput](publishPostNode.GetInputSchema())
	publishPostExecutor := publishpost.NewInstagramPublishPostExecutor(publishPostNodeID, publishPostValidator)

	nodes.RegisterNode(publishPostNode)
	executors.RegisterExecutor(publishPostExecutor)
	functioncalling.RegisterFunctionCalling(functioncalling.Generate(publishPostNode))

	publishStoryNodeID := nodeID + "_publish_story"
	publishStoryNode := publishstory.NewInstagramPublishStoryNode(publishStoryNodeID)
	publishStoryValidator := jsonschema.New[publishstory.InstagramPublishStoryExecutorInput](publishStoryNode.GetInputSchema())
	publishStoryExecutor := publishstory.NewInstagramPublishStoryExecutor(publishStoryNodeID, publishStoryValidator)

	nodes.RegisterNode(publishStoryNode)
	executors.RegisterExecutor(publishStoryExecutor)
	functioncalling.RegisterFunctionCalling(functioncalling.Generate(publishStoryNode))

	publishReelsNodeID := nodeID + "_publish_reels"
	publishReelsNode := publishreels.NewInstagramPublishReelsNode(publishReelsNodeID)
	publishReelsValidator := jsonschema.New[publishreels.InstagramPublishReelsExecutorInput](publishReelsNode.GetInputSchema())
	publishReelsExecutor := publishreels.NewInstagramPublishReelsExecutor(publishReelsNodeID, publishReelsValidator)

	nodes.RegisterNode(publishReelsNode)
	executors.RegisterExecutor(publishReelsExecutor)
	functioncalling.RegisterFunctionCalling(functioncalling.Generate(publishReelsNode))

	mcp.RegisterServer(&mcp.Server{
		ID:          nodeID,
		Name:        "Instagram",
		Description: "Tools for publishing posts, stories, and reels to Instagram.",
		Icon: mcp.ServerIcon{
			Brand: "instagram",
		},
		Version: "0.0.1",
		Tools: []nodes.NodeManager{
			publishPostNode,
			publishStoryNode,
			publishReelsNode,
		},
	})
}
