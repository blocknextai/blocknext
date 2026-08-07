package starter

import (
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
	gjs "github.com/google/jsonschema-go/jsonschema"
)

type StarterNode struct {
	nodes.Node
}

func NewStarterNode(nodeID string) *StarterNode {
	return &StarterNode{
		Node: nodes.Node{
			ID:          nodeID,
			Version:     "0.0.1",
			Name:        "Starter",
			Description: "Entry point that starts the workflow execution.",
			Icon: nodes.NodeIcon{
				Glyph: "play",
			},
			Inputs: []nodes.NodeHandle{},
			Outputs: []nodes.NodeHandle{
				{Key: "out"},
			},
			Categories:    []string{"System"},
			SubCategories: []string{"System"},
			Tags: []string{
				"flow",
				"control",
				"utility",
				"trigger",
				"start",
			},
			InputSchema: &gjs.Schema{
				Type:       "object",
				Properties: map[string]*gjs.Schema{},
			},
			OutputSchema: &gjs.Schema{
				Type: "array",
				Items: &gjs.Schema{
					Type: "object",
					Properties: map[string]*gjs.Schema{
						"status": {
							Type:        "boolean",
							Description: "Whether the starter executed successfully.",
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
