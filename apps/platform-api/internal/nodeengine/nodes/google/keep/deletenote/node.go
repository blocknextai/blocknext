package deletenote

import (
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
	gjs "github.com/google/jsonschema-go/jsonschema"
)

type GoogleKeepDeleteNoteNode struct {
	nodes.Node
}

func NewGoogleKeepDeleteNoteNode(nodeID string) *GoogleKeepDeleteNoteNode {
	return &GoogleKeepDeleteNoteNode{
		ID:          nodeID,
		Kind:        nodes.NodeKindAction,
		Version:     "0.0.1",
		Name:        "Google Keep Delete Note",
		Description: "Delete a note from Google Keep. This removes the note immediately and cannot be undone.",
		Icon: nodes.NodeIcon{
			Brand: "google_keep",
			Glyph: "trash",
		},
		Inputs: []nodes.NodeHandle{
			{Key: "in"},
		},
		Outputs: []nodes.NodeHandle{
			{Key: "out"},
		},
		Categories:    []string{"Google Workspace"},
		SubCategories: []string{"Google Keep"},
		Tags: []string{
			"google",
			"note",
			"delete",
			"remove",
			"notes",
		},
		SupportedCredentials: []string{
			"google_keep_oauth2",
		},
		InputSchema: &gjs.Schema{
			Type: "object",
			Properties: map[string]*gjs.Schema{
				"name": {
					Type:        "string",
					Title:       "Name",
					Description: "Resource name of the note to delete (e.g. notes/abc123).",
				},
			},
			Required: []string{
				"name",
			},
		},
		OutputSchema: &gjs.Schema{
			Type: "array",
			Items: &gjs.Schema{
				Type: "object",
				Properties: map[string]*gjs.Schema{
					"status": {
						Type:        "boolean",
						Description: "Whether the note was deleted successfully.",
					},
				},
			},
		},
		HasNaturalLanguage: true,
		Annotations: nodes.NodeAnnotations{
			Destructive: new(true),
		},
	}
}
