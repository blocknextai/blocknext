package readdocs

import (
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
	gjs "github.com/google/jsonschema-go/jsonschema"
)

type GoogleDocsReadDocsNode struct {
	nodes.Node
}

func NewGoogleDocsReadDocsNode(nodeID string) *GoogleDocsReadDocsNode {
	return &GoogleDocsReadDocsNode{
		Node: nodes.Node{
			ID:          nodeID,
			Kind:        nodes.NodeKindAction,
			Version:     "0.0.1",
			Name:        "Google Docs Read",
			Description: "Read the plain text content of a Google Docs document.",
			Icon: nodes.NodeIcon{
				Brand: "google_docs",
				Glyph: "eye",
			},
			Inputs: []nodes.NodeHandle{
				{Key: "in"},
			},
			Outputs: []nodes.NodeHandle{
				{Key: "out"},
			},
			Categories:    []string{"Google Workspace"},
			SubCategories: []string{"Google Docs"},
			Tags: []string{
				"google",
				"docs",
				"document",
				"read",
				"get",
				"productivity",
			},
			SupportedCredentials: []string{
				"google_docs_oauth2",
			},
			InputSchema: &gjs.Schema{
				Type: "object",
				Properties: map[string]*gjs.Schema{
					"documentId": {
						Type:        "string",
						Title:       "Document ID",
						Description: "Google Docs document identifier to read.",
					},
				},
				Required: []string{
					"documentId",
				},
			},
			OutputSchema: &gjs.Schema{
				Type: "array",
				Items: &gjs.Schema{
					Type: "object",
					Properties: map[string]*gjs.Schema{
						"text": {
							Type:        "string",
							Description: "Plain text content extracted from the document.",
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
