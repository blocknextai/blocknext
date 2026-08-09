package updatepage

import (
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
	gjs "github.com/google/jsonschema-go/jsonschema"
)

type NotionUpdatePageNode struct {
	nodes.Node
}

func NewNotionUpdatePageNode(nodeID string) *NotionUpdatePageNode {
	return &NotionUpdatePageNode{
		Node: nodes.Node{
			ID:          nodeID,
			Kind:        nodes.NodeKindAction,
			Version:     "0.0.1",
			Name:        "Notion Update Page",
			Description: "Update an existing page in a Notion workspace.",
			Icon: nodes.NodeIcon{
				Brand: "notion",
				Glyph: "pencil",
			},
			Inputs: []nodes.NodeHandle{
				{Key: "in"},
			},
			Outputs: []nodes.NodeHandle{
				{Key: "out"},
			},
			Categories:    []string{"Publishing"},
			SubCategories: []string{"Notion"},
			Tags: []string{
				"page",
				"update",
				"modify",
				"edit",
				"workspace",
				"productivity",
				"docs",
			},
			SupportedCredentials: []string{
				"notion_oauth2",
			},
			InputSchema: &gjs.Schema{
				Type: "object",
				Properties: map[string]*gjs.Schema{
					"pageId": {
						Type:        "string",
						Title:       "Page ID",
						Description: "Identifier of the Notion page to update.",
					},
					"title": {
						Type:        "string",
						Title:       "Title",
						Description: "New title to set on the Notion page.",
					},
					"properties": {
						Type:        "string",
						Title:       "Properties",
						Description: "JSON-encoded object of additional Notion properties to update.",
					},
				},
				Required: []string{
					"pageId",
				},
			},
			OutputSchema: &gjs.Schema{
				Type: "array",
				Items: &gjs.Schema{
					Type: "object",
					Properties: map[string]*gjs.Schema{
						"status": {
							Type:        "boolean",
							Description: "Whether the Notion page was updated successfully.",
						},
					},
				},
			},
			HasNaturalLanguage: true,
			Annotations: nodes.NodeAnnotations{
				Destructive: new(true),
			},
		},
	}
}
