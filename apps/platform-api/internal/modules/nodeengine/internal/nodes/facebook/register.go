package facebook

import (
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/application/jsonschema"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/mcp"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/providers"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/nodes/facebook/publishmediapost"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/nodes/facebook/publishpost"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/nodes/facebook/publishstory"
)

func Register() {
	nodeID := "facebook"

	publishPostNodeID := nodeID + "_publish_post"
	publishPostNode := publishpost.NewFacebookPublishPostNode(publishPostNodeID)
	publishPostValidator := jsonschema.New[publishpost.FacebookPublishPostExecutorInput](publishPostNode.GetInputSchema())
	publishPostExecutor := publishpost.NewFacebookPublishPostExecutor(publishPostNodeID, publishPostValidator)

	publishMediaPostNodeID := nodeID + "_publish_media_post"
	publishMediaPostNode := publishmediapost.NewFacebookPublishMediaPostNode(publishMediaPostNodeID)
	publishMediaPostValidator := jsonschema.New[publishmediapost.FacebookPublishMediaPostExecutorInput](publishMediaPostNode.GetInputSchema())
	publishMediaPostExecutor := publishmediapost.NewFacebookPublishMediaPostExecutor(publishMediaPostNodeID, publishMediaPostValidator)

	publishStoryNodeID := nodeID + "_publish_story"
	publishStoryNode := publishstory.NewFacebookPublishStoryNode(publishStoryNodeID)
	publishStoryValidator := jsonschema.New[publishstory.FacebookPublishStoryExecutorInput](publishStoryNode.GetInputSchema())
	publishStoryExecutor := publishstory.NewFacebookPublishStoryExecutor(publishStoryNodeID, publishStoryValidator)

	providers.Register(providers.Provider{
		Server: mcp.Server{
			ID:          nodeID,
			Name:        "Facebook",
			Description: "Tools for publishing posts, media, and stories to Facebook pages.",
			Icon: mcp.ServerIcon{
				Brand: "facebook",
			},
			Version: "0.0.1",
		},
		Actions: []providers.Action{
			{Node: publishPostNode, Executor: publishPostExecutor},
			{Node: publishMediaPostNode, Executor: publishMediaPostExecutor},
			{Node: publishStoryNode, Executor: publishStoryExecutor},
		},
	})
}
