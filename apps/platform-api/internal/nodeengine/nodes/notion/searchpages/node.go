package searchpages

import (
	"github.com/blocknextai/go-packages/json"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
	gjs "github.com/google/jsonschema-go/jsonschema"
)

type NotionSearchPagesNode struct {
	nodes.Node
}

func NewNotionSearchPagesNode(nodeID string) *NotionSearchPagesNode {
	return &NotionSearchPagesNode{
		ID:          nodeID,
		Kind:        nodes.NodeKindAction,
		Version:     "0.0.1",
		Name:        "Notion Search Pages",
		Description: "Search for pages in a Notion workspace.",
		Icon: nodes.NodeIcon{
			Brand: "notion",
			Glyph: "search",
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
			"search",
			"pages",
			"find",
			"query",
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
				"query": {
					Type:        "string",
					Title:       "Query",
					Description: "Search query used to filter Notion pages.",
				},
				"limit": {
					Type:        "number",
					Title:       "Limit",
					Description: "Maximum number of pages to return.",
					Default:     json.RawMessage(`10`),
				},
			},
			Required: []string{
				"query",
			},
		},
		OutputSchema: &gjs.Schema{
			Type: "array",
			Items: &gjs.Schema{
				Type: "object",
				Properties: map[string]*gjs.Schema{
					"results": {
						Type:        "array",
						Description: "Notion pages that matched the search query.",
						Items:       &gjs.Schema{Type: "object"},
					},
				},
			},
		},
		HasNaturalLanguage: true,
		Annotations: nodes.NodeAnnotations{
			ReadOnly:   true,
			Idempotent: true,
		},
	}
}
