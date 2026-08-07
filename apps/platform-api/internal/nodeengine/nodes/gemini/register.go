package gemini

import (
	"github.com/blocknextai/platform-api/internal/filegateway"
	"github.com/blocknextai/platform-api/internal/nodeengine/application/jsonschema"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/executors"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/functioncalling"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/mcp"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/gemini/chat"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/gemini/nanobanana"
)

func Register(fileGatewayService filegateway.FileGateway) {
	nodeID := "gemini"

	nanoBananaNodeID := nodeID + "_nano_banana"
	nanoBananaNode := nanobanana.NewGeminiImageGenerationNode(nanoBananaNodeID)
	nanoBananaValidator := jsonschema.New[nanobanana.GeminiImageGenerationExecutorInput](nanoBananaNode.GetInputSchema())
	nanoBananaExecutor := nanobanana.NewGeminiImageGenerationExecutor(nanoBananaNodeID, nanoBananaValidator, fileGatewayService)

	nodes.RegisterNode(nanoBananaNode)
	executors.RegisterExecutor(nanoBananaExecutor)
	functioncalling.RegisterFunctionCalling(functioncalling.Generate(nanoBananaNode))

	chatNodeID := nodeID + "_chat"
	chatNode := chat.NewGeminiChatNode(chatNodeID)
	chatValidator := jsonschema.New[chat.GeminiChatExecutorInput](chatNode.GetInputSchema())
	chatExecutor := chat.NewGeminiChatExecutor(chatNodeID, chatValidator)

	nodes.RegisterNode(chatNode)
	executors.RegisterExecutor(chatExecutor)
	functioncalling.RegisterFunctionCalling(functioncalling.Generate(chatNode))

	mcp.RegisterServer(&mcp.Server{
		ID:          nodeID,
		Name:        "Gemini",
		Description: "Tools for chat completions and image generation with Google Gemini.",
		Icon: mcp.ServerIcon{
			Brand: "gemini",
		},
		Version: "0.0.1",
		Tools: []nodes.NodeManager{
			nanoBananaNode,
			chatNode,
		},
	})
}
