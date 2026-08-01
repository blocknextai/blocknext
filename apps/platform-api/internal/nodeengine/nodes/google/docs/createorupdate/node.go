package createorupdate

import (
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
	gjs "github.com/google/jsonschema-go/jsonschema"
)

type GoogleDocsCreateOrUpdateNode struct {
	nodes.Node
}

func NewGoogleDocsCreateOrUpdateNode(nodeID string) *GoogleDocsCreateOrUpdateNode {
	return &GoogleDocsCreateOrUpdateNode{
		Node: nodes.Node{
			ID:          nodeID,
			Version:     "0.0.1",
			Name:        "Google Docs Create or Update",
			Description: "Create a new Google Docs document or update an existing one with the provided content.",
			Icon: nodes.NodeIcon{
				Light: nodeID,
				Dark:  nodeID,
			},
			Categories:    []string{"Google Workspace"},
			SubCategories: []string{"Google Docs"},
			Tags: []string{
				"google",
				"docs",
				"document",
				"create",
				"update",
				"productivity",
			},
			SupportedCredentials: []string{
				"google_docs_oauth2",
			},
			InputSchema: &gjs.Schema{
				Type: "object",
				Properties: map[string]*gjs.Schema{
					"title": {
						Type:        "string",
						Title:       "Title",
						Description: "Title to use when creating a new document.",
					},
					"content": {
						Type:        "string",
						Title:       "Content",
						Description: "Text content to insert into the document.",
					},
					"documentId": {
						Type:        "string",
						Title:       "Document ID",
						Description: "Existing Google Docs document identifier to update; leave empty to create a new document.",
					},
				},
				Required: []string{
					"content",
				},
			},
			OutputSchema: &gjs.Schema{
				Type: "array",
				Items: &gjs.Schema{
					Type: "object",
					Properties: map[string]*gjs.Schema{
						"status": {
							Type:        "boolean",
							Description: "Whether the document was created or updated successfully.",
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
