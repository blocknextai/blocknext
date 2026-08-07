package elevenlabs

import (
	"github.com/blocknextai/platform-api/internal/filegateway"
	"github.com/blocknextai/platform-api/internal/nodeengine/application/jsonschema"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/executors"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/functioncalling"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/mcp"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/elevenlabs/texttospeech"
)

func Register(fileGateway filegateway.FileGateway) {
	nodeID := "elevenlabs"

	textToSpeechNodeID := nodeID + "_text_to_speech"
	textToSpeechNode := texttospeech.NewElevenlabsTextToSpeechNode(textToSpeechNodeID)
	textToSpeechValidator := jsonschema.New[texttospeech.ElevenlabsTextToSpeechExecutorInput](textToSpeechNode.GetInputSchema())
	textToSpeechExecutor := texttospeech.NewElevenlabsTextToSpeechExecutor(textToSpeechNodeID, textToSpeechValidator, fileGateway)

	nodes.RegisterNode(textToSpeechNode)
	executors.RegisterExecutor(textToSpeechExecutor)
	functioncalling.RegisterFunctionCalling(functioncalling.Generate(textToSpeechNode))

	mcp.RegisterServer(&mcp.Server{
		ID:          nodeID,
		Name:        "ElevenLabs",
		Description: "Tools for text-to-speech and audio synthesis via ElevenLabs.",
		Icon: mcp.ServerIcon{
			Brand: "elevenlabs",
		},
		Version: "0.0.1",
		Tools: []nodes.NodeManager{
			textToSpeechNode,
		},
	})
}
