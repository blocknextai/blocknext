package providers

import (
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/adapters"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/executors"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/functioncalling"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/mcp"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/nodes"
)

type Action struct {
	Node     nodes.NodeManager
	Executor executors.ExecutorManager
}

type Provider struct {
	Server   mcp.Server
	Actions  []Action
	Adapters []adapters.TriggerAdapter
}

func Register(provider Provider) {
	tools := make([]nodes.NodeManager, 0, len(provider.Actions))
	for _, action := range provider.Actions {
		nodes.RegisterNode(action.Node)
		if action.Executor == nil {
			continue
		}

		executors.RegisterExecutor(action.Executor)
		if action.Node.GetHasNaturalLanguage() {
			functioncalling.RegisterFunctionCalling(functioncalling.Generate(action.Node))
		}
		tools = append(tools, action.Node)
	}

	for _, adapter := range provider.Adapters {
		adapters.RegisterAdapter(adapter)
	}

	server := provider.Server
	server.Tools = tools
	mcp.RegisterServer(&server)
}
