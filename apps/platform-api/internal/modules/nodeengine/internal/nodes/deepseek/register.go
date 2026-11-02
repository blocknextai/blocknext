package deepseek

import (
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/application/jsonschema"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/mcp"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/providers"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/nodes/deepseek/chat"
)

func Register() {
	nodeID := "deepseek"

	chatNodeID := nodeID + "_chat"
	chatNode := chat.NewDeepseekChatNode(chatNodeID)
	chatValidator := jsonschema.New[chat.DeepseekChatExecutorInput](chatNode.GetInputSchema())
	chatExecutor := chat.NewDeepseekChatExecutor(chatNodeID, chatValidator)

	providers.Register(providers.Provider{
		Server: mcp.Server{
			ID:          nodeID,
			Name:        "DeepSeek",
			Description: "Tools for chat completions with DeepSeek models.",
			Icon: mcp.ServerIcon{
				Brand: "deepseek",
			},
			Version: "0.0.1",
		},
		Actions: []providers.Action{
			{Node: chatNode, Executor: chatExecutor},
		},
	})
}
