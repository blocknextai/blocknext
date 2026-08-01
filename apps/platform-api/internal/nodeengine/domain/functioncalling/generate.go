package functioncalling

import (
	"strings"

	"github.com/blocknextai/go-packages/json"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
)

func Generate(node nodes.NodeManager) *FunctionCalling {
	return &FunctionCalling{
		ID:       node.GetID(),
		Disabled: !node.GetHasNaturalLanguage(),
		FunctionDeclarations: []map[string]any{
			fromNode(node),
		},
	}
}

func fromNode(node nodes.NodeManager) map[string]any {
	name := strings.ReplaceAll(node.GetID(), ".", "_")

	parameters := map[string]any{
		"type":       "object",
		"properties": map[string]any{},
	}
	if input := node.GetInputSchema(); input != nil {
		if data, err := json.Marshal(input); err == nil {
			var asMap map[string]any
			if err := json.Unmarshal(data, &asMap); err == nil {
				parameters = asMap
			}
		}
	}

	return map[string]any{
		"name":        name,
		"description": node.GetDescription(),
		"parameters":  parameters,
	}
}
