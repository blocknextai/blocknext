package chatgpt

import (
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/application/jsonschema"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/mcp"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/providers"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/nodes/chatgpt/chat"
)

func Register() {
	nodeID := "chatgpt"

	chatNodeID := nodeID + "_chat"
	chatNode := chat.NewChatgptChatNode(chatNodeID)
	chatValidator := jsonschema.New[chat.ChatgptChatExecutorInput](chatNode.GetInputSchema())
	chatExecutor := chat.NewChatgptChatExecutor(chatNodeID, chatValidator)

	providers.Register(providers.Provider{
		Server: mcp.Server{
			ID:          nodeID,
			Name:        "ChatGPT",
			Description: "Tools for chat completions with OpenAI ChatGPT.",
			Icon: mcp.ServerIcon{
				Brand: "chatgpt",
			},
			Version: "0.0.1",
		},
		Actions: []providers.Action{
			{Node: chatNode, Executor: chatExecutor},
		},
	})
}
