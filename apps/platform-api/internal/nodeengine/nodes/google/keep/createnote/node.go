package createnote

import (
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/google/keep/helpers"
	gjs "github.com/google/jsonschema-go/jsonschema"
)

type GoogleKeepCreateNoteNode struct {
	nodes.Node
}

func NewGoogleKeepCreateNoteNode(nodeID string) *GoogleKeepCreateNoteNode {
	return &GoogleKeepCreateNoteNode{
		Node: nodes.Node{
			ID:          nodeID,
			Kind:        nodes.NodeKindAction,
			Version:     "0.0.1",
			Name:        "Google Keep Create Note",
			Description: "Create a new note in Google Keep.",
			Icon: nodes.NodeIcon{
				Brand: "google_keep",
				Glyph: "note",
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
				"create",
				"productivity",
				"notes",
			},
			SupportedCredentials: []string{
				"google_keep_oauth2",
			},
			InputSchema: &gjs.Schema{
				Type: "object",
				Properties: map[string]*gjs.Schema{
					"title": {
						Type:        "string",
						Title:       "Title",
						Description: "Title of the note.",
					},
					"content": {
						Type:        "string",
						Title:       "Content",
						Description: "Body content of the note.",
					},
				},
				Required: []string{
					"title",
					"content",
				},
			},
			OutputSchema: &gjs.Schema{
				Type:  "array",
				Items: helpers.NoteSchema(),
			},
			HasNaturalLanguage: true,
			Annotations: nodes.NodeAnnotations{
				Destructive: new(false),
			},
		},
	}
}
