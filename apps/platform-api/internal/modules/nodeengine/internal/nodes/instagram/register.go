package instagram

import (
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/application/jsonschema"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/mcp"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/providers"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/nodes/instagram/publishpost"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/nodes/instagram/publishreels"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/nodes/instagram/publishstory"
)

func Register() {
	nodeID := "instagram"

	publishPostNodeID := nodeID + "_publish_post"
	publishPostNode := publishpost.NewInstagramPublishPostNode(publishPostNodeID)
	publishPostValidator := jsonschema.New[publishpost.InstagramPublishPostExecutorInput](publishPostNode.GetInputSchema())
	publishPostExecutor := publishpost.NewInstagramPublishPostExecutor(publishPostNodeID, publishPostValidator)

	publishStoryNodeID := nodeID + "_publish_story"
	publishStoryNode := publishstory.NewInstagramPublishStoryNode(publishStoryNodeID)
	publishStoryValidator := jsonschema.New[publishstory.InstagramPublishStoryExecutorInput](publishStoryNode.GetInputSchema())
	publishStoryExecutor := publishstory.NewInstagramPublishStoryExecutor(publishStoryNodeID, publishStoryValidator)

	publishReelsNodeID := nodeID + "_publish_reels"
	publishReelsNode := publishreels.NewInstagramPublishReelsNode(publishReelsNodeID)
	publishReelsValidator := jsonschema.New[publishreels.InstagramPublishReelsExecutorInput](publishReelsNode.GetInputSchema())
	publishReelsExecutor := publishreels.NewInstagramPublishReelsExecutor(publishReelsNodeID, publishReelsValidator)

	providers.Register(providers.Provider{
		Server: mcp.Server{
			ID:          nodeID,
			Name:        "Instagram",
			Description: "Tools for publishing posts, stories, and reels to Instagram.",
			Icon: mcp.ServerIcon{
				Brand: "instagram",
			},
			Version: "0.0.1",
		},
		Actions: []providers.Action{
			{Node: publishPostNode, Executor: publishPostExecutor},
			{Node: publishStoryNode, Executor: publishStoryExecutor},
			{Node: publishReelsNode, Executor: publishReelsExecutor},
		},
	})
}
