package anthropic

import (
	"github.com/blocknextai/platform-api/internal/nodeengine/application/jsonschema"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/executors"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/functioncalling"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/mcp"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/anthropic/chat"
)

func Register() {
	nodeID := "anthropic"

	chatNodeID := nodeID + "_chat"
	chatNode := chat.NewAnthropicChatNode(chatNodeID)
	chatValidator := jsonschema.New[chat.AnthropicChatExecutorInput](chatNode.GetInputSchema())
	chatExecutor := chat.NewAnthropicChatExecutor(chatNodeID, chatValidator)

	nodes.RegisterNode(chatNode)
	executors.RegisterExecutor(chatExecutor)
	functioncalling.RegisterFunctionCalling(functioncalling.Generate(chatNode))

	mcp.RegisterServer(&mcp.Server{
		ID:          nodeID,
		Name:        "Anthropic",
		Description: "Tools for chat completions with Anthropic Claude.",
		Version:     "0.0.1",
		Tools: []nodes.NodeManager{
			chatNode,
		},
	})
}
