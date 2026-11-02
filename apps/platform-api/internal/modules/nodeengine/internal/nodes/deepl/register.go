package deepl

import (
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/application/jsonschema"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/mcp"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/providers"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/nodes/deepl/translate"
)

func Register() {
	nodeID := "deepl"

	translateNodeID := nodeID + "_translate"
	translateNode := translate.NewDeeplTranslateNode(translateNodeID)
	translateValidator := jsonschema.New[translate.DeeplTranslateExecutorInput](translateNode.GetInputSchema())
	translateExecutor := translate.NewDeeplTranslateExecutor(translateNodeID, translateValidator)

	providers.Register(providers.Provider{
		Server: mcp.Server{
			ID:          nodeID,
			Name:        "DeepL",
			Description: "Tools for translation via DeepL.",
			Icon: mcp.ServerIcon{
				Brand: "deepl",
			},
			Version: "0.0.1",
		},
		Actions: []providers.Action{
			{Node: translateNode, Executor: translateExecutor},
		},
	})
}
