package youtube

import (
	"github.com/blocknextai/platform-api/internal/filegateway"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/application/jsonschema"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/mcp"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/providers"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/nodes/google/youtube/uploadvideo"
)

func Register(fileGateway filegateway.FileGateway) {
	nodeID := "youtube"

	uploadVideoNodeID := nodeID + "_upload_video"
	uploadVideoNode := uploadvideo.NewYouTubeUploadVideoNode(uploadVideoNodeID)
	uploadVideoValidator := jsonschema.New[uploadvideo.YouTubeUploadVideoExecutorInput](uploadVideoNode.GetInputSchema())
	uploadVideoExecutor := uploadvideo.NewYouTubeUploadVideoExecutor(uploadVideoNodeID, uploadVideoValidator, fileGateway)

	providers.Register(providers.Provider{
		Server: mcp.Server{
			ID:          nodeID,
			Name:        "YouTube",
			Description: "Tools for uploading videos to YouTube.",
			Icon: mcp.ServerIcon{
				Brand: "youtube",
			},
			Version: "0.0.1",
		},
		Actions: []providers.Action{
			{Node: uploadVideoNode, Executor: uploadVideoExecutor},
		},
	})
}
