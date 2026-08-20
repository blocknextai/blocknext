package getnote

import (
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/google/keep/helpers"
	gjs "github.com/google/jsonschema-go/jsonschema"
)

type GoogleKeepGetNoteNode struct {
	nodes.Node
}

func NewGoogleKeepGetNoteNode(nodeID string) *GoogleKeepGetNoteNode {
	return &GoogleKeepGetNoteNode{
		ID:          nodeID,
		Kind:        nodes.NodeKindAction,
		Version:     "0.0.1",
		Name:        "Google Keep Get Note",
		Description: "Get a single note from Google Keep by its resource name.",
		Icon: nodes.NodeIcon{
			Brand: "google_keep",
			Glyph: "eye",
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
			"get",
			"read",
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
					Description: "Resource name of the note to get (e.g. notes/abc123).",
				},
			},
			Required: []string{
				"name",
			},
		},
		OutputSchema: &gjs.Schema{
			Type:  "array",
			Items: helpers.NoteSchema(),
		},
		HasNaturalLanguage: true,
		Annotations: nodes.NodeAnnotations{
			ReadOnly:   true,
			Idempotent: true,
		},
	}
}
