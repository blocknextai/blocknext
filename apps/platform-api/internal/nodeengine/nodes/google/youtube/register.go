package youtube

import (
	"github.com/blocknextai/platform-api/internal/filegateway"
	"github.com/blocknextai/platform-api/internal/nodeengine/application/jsonschema"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/executors"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/functioncalling"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/mcp"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/google/youtube/uploadvideo"
)

func Register(fileGateway filegateway.FileGateway) {
	nodeID := "youtube"

	uploadVideoNodeID := nodeID + "_upload_video"
	uploadVideoNode := uploadvideo.NewYouTubeUploadVideoNode(uploadVideoNodeID)
	uploadVideoValidator := jsonschema.New[uploadvideo.YouTubeUploadVideoExecutorInput](uploadVideoNode.GetInputSchema())
	uploadVideoExecutor := uploadvideo.NewYouTubeUploadVideoExecutor(uploadVideoNodeID, uploadVideoValidator, fileGateway)

	nodes.RegisterNode(uploadVideoNode)
	executors.RegisterExecutor(uploadVideoExecutor)
	functioncalling.RegisterFunctionCalling(functioncalling.Generate(uploadVideoNode))

	mcp.RegisterServer(&mcp.Server{
		ID:          nodeID,
		Name:        "YouTube",
		Description: "Tools for uploading videos to YouTube.",
		Version:     "0.0.1",
		Tools: []nodes.NodeManager{
			uploadVideoNode,
		},
	})
}
