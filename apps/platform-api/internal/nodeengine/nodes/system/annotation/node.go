package annotation

import (
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
	gjs "github.com/google/jsonschema-go/jsonschema"
)

type AnnotationNode struct {
	nodes.Node
}

func NewAnnotationNode(nodeID string) *AnnotationNode {
	return &AnnotationNode{
		ID:          nodeID,
		Kind:        nodes.NodeKindNote,
		Version:     "0.0.1",
		Name:        "Annotation",
		Description: "A free-text note pinned to the canvas.",
		Icon: nodes.NodeIcon{
			Glyph: "note",
		},
		Inputs:        []nodes.NodeHandle{},
		Outputs:       []nodes.NodeHandle{},
		Categories:    []string{"System"},
		SubCategories: []string{"System"},
		Tags: []string{
			"note",
			"comment",
			"documentation",
		},
		InputSchema: &gjs.Schema{
			Type: "object",
			Properties: map[string]*gjs.Schema{
				"note": {
					Type:        "string",
					Description: "The text shown on the canvas.",
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
	}
}
