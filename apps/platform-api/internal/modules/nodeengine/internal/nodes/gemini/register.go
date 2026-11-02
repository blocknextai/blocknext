package gemini

import (
	"github.com/blocknextai/platform-api/internal/filegateway"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/application/jsonschema"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/mcp"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/providers"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/nodes/gemini/chat"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/nodes/gemini/nanobanana"
)

func Register(fileGatewayService filegateway.FileGateway) {
	nodeID := "gemini"

	nanoBananaNodeID := nodeID + "_nano_banana"
	nanoBananaNode := nanobanana.NewGeminiImageGenerationNode(nanoBananaNodeID)
	nanoBananaValidator := jsonschema.New[nanobanana.GeminiImageGenerationExecutorInput](nanoBananaNode.GetInputSchema())
	nanoBananaExecutor := nanobanana.NewGeminiImageGenerationExecutor(nanoBananaNodeID, nanoBananaValidator, fileGatewayService)

	chatNodeID := nodeID + "_chat"
	chatNode := chat.NewGeminiChatNode(chatNodeID)
	chatValidator := jsonschema.New[chat.GeminiChatExecutorInput](chatNode.GetInputSchema())
	chatExecutor := chat.NewGeminiChatExecutor(chatNodeID, chatValidator)

	providers.Register(providers.Provider{
		Server: mcp.Server{
			ID:          nodeID,
			Name:        "Gemini",
			Description: "Tools for chat completions and image generation with Google Gemini.",
			Icon: mcp.ServerIcon{
				Brand: "gemini",
			},
			Version: "0.0.1",
		},
		Actions: []providers.Action{
			{Node: nanoBananaNode, Executor: nanoBananaExecutor},
			{Node: chatNode, Executor: chatExecutor},
		},
	})
}
