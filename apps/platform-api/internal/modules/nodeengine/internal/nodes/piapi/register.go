package piapi

import (
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/application/jsonschema"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/executors"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/functioncalling"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/mcp"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/nodes"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/nodes/piapi/audiogen"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/nodes/piapi/imagegen"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/nodes/piapi/videogen"
)

func Register() {
	nodeID := "piapi"

	imageGenNodeID := nodeID + "_image_gen"
	imageGenNode := imagegen.NewPiAPIImageGenNode(imageGenNodeID)
	imageGenValidator := jsonschema.New[imagegen.PiAPIImageGenInput](imageGenNode.GetInputSchema())
	imageGenExecutor := imagegen.NewPiAPIImageGenExecutor(imageGenNodeID, imageGenValidator)

	nodes.RegisterNode(imageGenNode)
	executors.RegisterExecutor(imageGenExecutor)
	functioncalling.RegisterFunctionCalling(functioncalling.Generate(imageGenNode))

	videoGenNodeID := nodeID + "_video_gen"
	videoGenNode := videogen.NewPiAPIVideoGenNode(videoGenNodeID)
	videoGenValidator := jsonschema.New[videogen.PiAPIVideoGenInput](videoGenNode.GetInputSchema())
	videoGenExecutor := videogen.NewPiAPIVideoGenExecutor(videoGenNodeID, videoGenValidator)

	nodes.RegisterNode(videoGenNode)
	executors.RegisterExecutor(videoGenExecutor)
	functioncalling.RegisterFunctionCalling(functioncalling.Generate(videoGenNode))

	audioGenNodeID := nodeID + "_audio_gen"
	audioGenNode := audiogen.NewPiAPIAudioGenNode(audioGenNodeID)
	audioGenValidator := jsonschema.New[audiogen.PiAPIAudioGenInput](audioGenNode.GetInputSchema())
	audioGenExecutor := audiogen.NewPiAPIAudioGenExecutor(audioGenNodeID, audioGenValidator)

	nodes.RegisterNode(audioGenNode)
	executors.RegisterExecutor(audioGenExecutor)
	functioncalling.RegisterFunctionCalling(functioncalling.Generate(audioGenNode))

	mcp.RegisterServer(&mcp.Server{
		ID:          nodeID,
		Name:        "PiAPI",
		Description: "Tools for AI image, video, and audio generation via PiAPI.",
		Icon: mcp.ServerIcon{
			Brand: "piapi",
		},
		Version: "0.0.1",
		Tools: []nodes.NodeManager{
			imageGenNode,
			videoGenNode,
			audioGenNode,
		},
	})
}
