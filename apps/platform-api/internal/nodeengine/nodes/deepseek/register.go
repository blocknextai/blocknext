package deepseek

import (
	"github.com/blocknextai/platform-api/internal/nodeengine/application/jsonschema"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/executors"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/functioncalling"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/mcp"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/deepseek/chat"
)

func Register() {
	nodeID := "deepseek"

	chatNodeID := nodeID + "_chat"
	chatNode := chat.NewDeepseekChatNode(chatNodeID)
	chatValidator := jsonschema.New[chat.DeepseekChatExecutorInput](chatNode.GetInputSchema())
	chatExecutor := chat.NewDeepseekChatExecutor(chatNodeID, chatValidator)

	nodes.RegisterNode(chatNode)
	executors.RegisterExecutor(chatExecutor)
	functioncalling.RegisterFunctionCalling(functioncalling.Generate(chatNode))

	mcp.RegisterServer(&mcp.Server{
		ID:          nodeID,
		Name:        "DeepSeek",
		Description: "Tools for chat completions with DeepSeek models.",
		Icon: mcp.ServerIcon{
			Brand: "deepseek",
		},
		Version: "0.0.1",
		Tools: []nodes.NodeManager{
			chatNode,
		},
	})
}
