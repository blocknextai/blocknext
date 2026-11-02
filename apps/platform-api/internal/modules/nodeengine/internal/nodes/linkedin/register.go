package linkedin

import (
	"github.com/blocknextai/platform-api/internal/filegateway"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/application/jsonschema"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/mcp"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/providers"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/nodes/linkedin/publishmediapost"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/nodes/linkedin/publishpost"
)

func Register(fileGateway filegateway.FileGateway) {
	nodeID := "linkedin"

	publishPostNodeID := nodeID + "_publish_post"
	publishPostNode := publishpost.NewLinkedinPublishPostNode(publishPostNodeID)
	publishPostValidator := jsonschema.New[publishpost.LinkedinPublishPostExecutorInput](publishPostNode.GetInputSchema())
	publishPostExecutor := publishpost.NewLinkedinPublishPostExecutor(publishPostNodeID, publishPostValidator)

	publishMediaPostNodeID := nodeID + "_publish_media_post"
	publishMediaPostNode := publishmediapost.NewLinkedinPublishMediaPostNode(publishMediaPostNodeID)
	publishMediaPostValidator := jsonschema.New[publishmediapost.LinkedinPublishMediaPostExecutorInput](publishMediaPostNode.GetInputSchema())
	publishMediaPostExecutor := publishmediapost.NewLinkedinPublishMediaPostExecutor(publishMediaPostNodeID, publishMediaPostValidator, fileGateway)

	providers.Register(providers.Provider{
		Server: mcp.Server{
			ID:          nodeID,
			Name:        "LinkedIn",
			Description: "Tools for publishing posts and media to LinkedIn.",
			Icon: mcp.ServerIcon{
				Brand: "linkedin",
			},
			Version: "0.0.1",
		},
		Actions: []providers.Action{
			{Node: publishPostNode, Executor: publishPostExecutor},
			{Node: publishMediaPostNode, Executor: publishMediaPostExecutor},
		},
	})
}
