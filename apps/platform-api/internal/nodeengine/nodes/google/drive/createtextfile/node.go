package createtextfile

import (
	"github.com/blocknextai/go-packages/json"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
	gjs "github.com/google/jsonschema-go/jsonschema"
)

type GoogleDriveCreateTextFileNode struct {
	nodes.Node
}

func NewGoogleDriveCreateTextFileNode(nodeID string) *GoogleDriveCreateTextFileNode {
	return &GoogleDriveCreateTextFileNode{
		Node: nodes.Node{
			ID:          nodeID,
			Version:     "0.0.1",
			Name:        "Google Drive Create Text File",
			Description: "Create a plain text file in Google Drive.",
			Icon: nodes.NodeIcon{
				Brand: "google_drive",
				Glyph: "plus",
			},
			Inputs: []nodes.NodeHandle{
				{Key: "in"},
			},
			Outputs: []nodes.NodeHandle{
				{Key: "out"},
			},
			Categories:    []string{"Google Workspace"},
			SubCategories: []string{"Google Drive"},
			Tags: []string{
				"google",
				"drive",
				"file",
				"create",
				"text",
				"storage",
			},
			SupportedCredentials: []string{
				"google_drive_oauth2",
			},
			InputSchema: &gjs.Schema{
				Type: "object",
				Properties: map[string]*gjs.Schema{
					"fileName": {
						Type:        "string",
						Title:       "File Name",
						Description: "Name of the text file to create.",
					},
					"content": {
						Type:        "string",
						Title:       "Content",
						Description: "Text content to write to the file.",
					},
					"parentFolderId": {
						Type:        "string",
						Title:       "Parent Folder ID",
						Description: "Identifier of the parent folder; use \"root\" for the drive root.",
						Default:     json.RawMessage(`"root"`),
					},
				},
				Required: []string{
					"fileName",
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
							Description: "Whether the file was created successfully.",
						},
						"fileId": {
							Type:        "string",
							Description: "Identifier of the created file.",
						},
						"fileName": {
							Type:        "string",
							Description: "Name of the created file.",
						},
						"webViewLink": {
							Type:        "string",
							Format:      "uri",
							Description: "Web view link to open the file in Google Drive.",
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
