package getpage

import (
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
	gjs "github.com/google/jsonschema-go/jsonschema"
)

type NotionGetPageNode struct {
	nodes.Node
}

func NewNotionGetPageNode(nodeID string) *NotionGetPageNode {
	return &NotionGetPageNode{
		Node: nodes.Node{
			ID:          nodeID,
			Version:     "0.0.1",
			Name:        "Notion Get Page",
			Description: "Retrieve a page from a Notion workspace.",
			Icon: nodes.NodeIcon{
				Light: nodeID,
				Dark:  nodeID,
			},
			Categories:    []string{"Publishing"},
			SubCategories: []string{"Notion"},
			Tags: []string{
				"page",
				"get",
				"retrieve",
				"read",
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
						Description: "Identifier of the Notion page to retrieve.",
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
						"id": {
							Type:        "string",
							Description: "Identifier of the Notion page.",
						},
						"object": {
							Type:        "string",
							Description: "Object type of the resource.",
						},
						"url": {
							Type:        "string",
							Description: "URL of the Notion page.",
							Format:      "uri",
						},
						"properties": {
							Type:        "object",
							Description: "Notion page properties keyed by name.",
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
