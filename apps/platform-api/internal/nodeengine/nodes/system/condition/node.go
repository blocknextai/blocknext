package condition

import (
	"github.com/blocknextai/go-packages/json"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
	gjs "github.com/google/jsonschema-go/jsonschema"
)

type ConditionNode struct {
	nodes.Node
}

func NewConditionNode(nodeID string) *ConditionNode {
	return &ConditionNode{
		Node: nodes.Node{
			ID:          nodeID,
			Kind:        nodes.NodeKindAction,
			Version:     "0.0.1",
			Name:        "Condition",
			Description: "Branch the workflow based on a conditional expression.",
			Icon: nodes.NodeIcon{
				Glyph: "branch",
			},
			Inputs: []nodes.NodeHandle{
				{Key: "in"},
			},
			Outputs: []nodes.NodeHandle{
				{Key: "true", Label: "True"},
				{Key: "false", Label: "False"},
			},
			Categories:    []string{"System"},
			SubCategories: []string{"System"},
			Tags: []string{
				"flow",
				"control",
				"utility",
				"branch",
				"if",
			},
			InputSchema: &gjs.Schema{
				Type: "object",
				Properties: map[string]*gjs.Schema{
					"leftValue": {
						Type:        "string",
						Title:       "Left Value",
						Description: "Left-hand operand of the condition.",
					},
					"operator": {
						Type:        "string",
						Title:       "Operator",
						Description: "Comparison operator to apply between the operands.",
						Enum: []any{
							"eq",
							"neq",
							"gt",
							"gte",
							"lt",
							"lte",
							"contains",
							"not_contains",
							"is_empty",
							"is_not_empty",
						},
						Default: json.RawMessage(`"eq"`),
					},
					"rightValue": {
						Type:        "string",
						Title:       "Right Value",
						Description: "Right-hand operand of the condition.",
					},
				},
				Required: []string{
					"operator",
				},
			},
			OutputSchema: &gjs.Schema{
				Type: "array",
				Items: &gjs.Schema{
					Type: "object",
					Properties: map[string]*gjs.Schema{
						"status": {
							Type:        "boolean",
							Description: "Boolean result of the condition evaluation.",
						},
					},
				},
			},
			HasNaturalLanguage: false,
			Annotations: nodes.NodeAnnotations{
				ReadOnly:    true,
				Idempotent:  true,
				Destructive: new(false),
				OpenWorld:   new(false),
			},
		},
	}
}
