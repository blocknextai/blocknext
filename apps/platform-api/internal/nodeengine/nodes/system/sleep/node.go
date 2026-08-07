package sleep

import (
	"github.com/blocknextai/go-packages/json"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
	gjs "github.com/google/jsonschema-go/jsonschema"
)

type SleepNode struct {
	nodes.Node
}

func NewSleepNode(nodeID string) *SleepNode {
	return &SleepNode{
		Node: nodes.Node{
			ID:          nodeID,
			Version:     "0.0.1",
			Name:        "Sleep",
			Description: "Pause the workflow for a specified duration.",
			Icon: nodes.NodeIcon{
				Glyph: "clock",
			},
			Inputs: []nodes.NodeHandle{
				{Key: "in"},
			},
			Outputs: []nodes.NodeHandle{
				{Key: "out"},
			},
			Categories:    []string{"System"},
			SubCategories: []string{"System"},
			Tags: []string{
				"flow",
				"control",
				"utility",
				"delay",
				"wait",
			},
			InputSchema: &gjs.Schema{
				Type: "object",
				Properties: map[string]*gjs.Schema{
					"duration": {
						Type:        "string",
						Title:       "Duration",
						Description: "Sleep duration in milliseconds.",
						Enum:        []any{"1000", "5000", "10000", "30000", "60000"},
						Default:     json.RawMessage(`"1000"`),
					},
				},
				Required: []string{
					"duration",
				},
			},
			OutputSchema: &gjs.Schema{
				Type: "array",
				Items: &gjs.Schema{
					Type: "object",
					Properties: map[string]*gjs.Schema{
						"status": {
							Type:        "boolean",
							Description: "Whether the sleep completed successfully.",
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
