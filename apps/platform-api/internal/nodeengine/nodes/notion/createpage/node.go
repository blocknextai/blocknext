package createpage

import (
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
	gjs "github.com/google/jsonschema-go/jsonschema"
)

type NotionCreatePageNode struct {
	nodes.Node
}

func NewNotionCreatePageNode(nodeID string) *NotionCreatePageNode {
	return &NotionCreatePageNode{
		Node: nodes.Node{
			ID:          nodeID,
			Kind:        nodes.NodeKindAction,
			Version:     "0.0.1",
			Name:        "Notion Create Page",
			Description: "Create a new page in a Notion workspace.",
			Icon: nodes.NodeIcon{
				Brand: "notion",
				Glyph: "plus",
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
				"create",
				"documentation",
				"workspace",
				"productivity",
				"docs",
				"notes",
			},
			SupportedCredentials: []string{
				"notion_oauth2",
			},
			InputSchema: &gjs.Schema{
				Type: "object",
				Properties: map[string]*gjs.Schema{
					"parentId": {
						Type:        "string",
						Title:       "Parent ID",
						Description: "Identifier of the parent Notion page.",
					},
					"title": {
						Type:        "string",
						Title:       "Title",
						Description: "Title of the new Notion page.",
					},
					"content": {
						Type:        "string",
						Title:       "Content",
						Description: "Optional initial body text for the page.",
					},
				},
				Required: []string{
					"parentId",
					"title",
				},
			},
			OutputSchema: &gjs.Schema{
				Type: "array",
				Items: &gjs.Schema{
					Type: "object",
					Properties: map[string]*gjs.Schema{
						"id": {
							Type:        "string",
							Description: "Identifier of the created Notion page.",
						},
					},
				},
			},
			HasNaturalLanguage: true,
			Annotations: nodes.NodeAnnotations{
				Destructive: new(false),
			},
		},
	}
}
