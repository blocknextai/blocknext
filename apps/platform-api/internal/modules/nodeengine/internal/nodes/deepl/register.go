package deepl

import (
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/application/jsonschema"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/executors"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/functioncalling"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/mcp"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/nodes"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/nodes/deepl/translate"
)

func Register() {
	nodeID := "deepl"

	translateNodeID := nodeID + "_translate"
	translateNode := translate.NewDeeplTranslateNode(translateNodeID)
	translateValidator := jsonschema.New[translate.DeeplTranslateExecutorInput](translateNode.GetInputSchema())
	translateExecutor := translate.NewDeeplTranslateExecutor(translateNodeID, translateValidator)

	nodes.RegisterNode(translateNode)
	executors.RegisterExecutor(translateExecutor)
	functioncalling.RegisterFunctionCalling(functioncalling.Generate(translateNode))

	mcp.RegisterServer(&mcp.Server{
		ID:          nodeID,
		Name:        "DeepL",
		Description: "Tools for translation via DeepL.",
		Icon: mcp.ServerIcon{
			Brand: "deepl",
		},
		Version: "0.0.1",
		Tools: []nodes.NodeManager{
			translateNode,
		},
	})
}
