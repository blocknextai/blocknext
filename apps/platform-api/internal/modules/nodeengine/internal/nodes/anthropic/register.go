package anthropic

import (
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/application/jsonschema"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/mcp"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/providers"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/nodes/anthropic/chat"
)

func Register() {
	nodeID := "anthropic"

	chatNodeID := nodeID + "_chat"
	chatNode := chat.NewAnthropicChatNode(chatNodeID)
	chatValidator := jsonschema.New[chat.AnthropicChatExecutorInput](chatNode.GetInputSchema())
	chatExecutor := chat.NewAnthropicChatExecutor(chatNodeID, chatValidator)

	providers.Register(providers.Provider{
		Server: mcp.Server{
			ID:          nodeID,
			Name:        "Anthropic",
			Description: "Tools for chat completions with Anthropic Claude.",
			Icon: mcp.ServerIcon{
				Brand: "anthropic",
			},
			Version: "0.0.1",
		},
		Actions: []providers.Action{
			{Node: chatNode, Executor: chatExecutor},
		},
	})
}
