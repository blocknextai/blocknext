package chatgpt

import (
	"github.com/blocknextai/platform-api/internal/nodeengine/application/jsonschema"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/executors"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/functioncalling"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/mcp"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/chatgpt/chat"
)

func Register() {
	nodeID := "chatgpt"

	chatNodeID := nodeID + "_chat"
	chatNode := chat.NewChatgptChatNode(chatNodeID)
	chatValidator := jsonschema.New[chat.ChatgptChatExecutorInput](chatNode.GetInputSchema())
	chatExecutor := chat.NewChatgptChatExecutor(chatNodeID, chatValidator)

	nodes.RegisterNode(chatNode)
	executors.RegisterExecutor(chatExecutor)
	functioncalling.RegisterFunctionCalling(functioncalling.Generate(chatNode))

	mcp.RegisterServer(&mcp.Server{
		ID:          nodeID,
		Name:        "ChatGPT",
		Description: "Tools for chat completions with OpenAI ChatGPT.",
		Version:     "0.0.1",
		Tools: []nodes.NodeManager{
			chatNode,
		},
	})
}
