package listnotes

import (
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/google/keep/helpers"
	gjs "github.com/google/jsonschema-go/jsonschema"
)

type GoogleKeepListNotesNode struct {
	nodes.Node
}

func NewGoogleKeepListNotesNode(nodeID string) *GoogleKeepListNotesNode {
	return &GoogleKeepListNotesNode{
		Node: nodes.Node{
			ID:          nodeID,
			Version:     "0.0.1",
			Name:        "Google Keep List Notes",
			Description: "List notes from Google Keep with optional filtering and pagination.",
			Icon: nodes.NodeIcon{
				Brand: "google_keep",
				Glyph: "list",
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
				"list",
				"read",
				"notes",
			},
			SupportedCredentials: []string{
				"google_keep_oauth2",
			},
			InputSchema: &gjs.Schema{
				Type: "object",
				Properties: map[string]*gjs.Schema{
					"pageSize": {
						Type:        "number",
						Title:       "Page Size",
						Description: "Maximum number of notes to return. A value of zero lets the server choose the upper bound.",
					},
					"pageToken": {
						Type:        "string",
						Title:       "Page Token",
						Description: "The previous page's nextPageToken, used to fetch the next page of results.",
					},
					"filter": {
						Type:        "string",
						Title:       "Filter",
						Description: "Filter for list results following the Google AIP filtering spec. Valid fields: createTime, updateTime, trashTime, trashed. Defaults to non-trashed notes.",
					},
				},
			},
			OutputSchema: &gjs.Schema{
				Type: "array",
				Items: &gjs.Schema{
					Type: "object",
					Properties: map[string]*gjs.Schema{
						"notes": {
							Type:        "array",
							Description: "A page of notes.",
							Items:       helpers.NoteSchema(),
						},
						"nextPageToken": {
							Type:        "string",
							Description: "Token to fetch the next page of results, empty when there are no more notes.",
						},
					},
				},
			},
			HasNaturalLanguage: true,
			Annotations: nodes.NodeAnnotations{
				ReadOnly:   true,
				Idempotent: true,
			},
		},
	}
}
