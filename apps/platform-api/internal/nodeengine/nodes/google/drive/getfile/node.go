package getfile

import (
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
	gjs "github.com/google/jsonschema-go/jsonschema"
)

type GoogleDriveGetFileNode struct {
	nodes.Node
}

func NewGoogleDriveGetFileNode(nodeID string) *GoogleDriveGetFileNode {
	return &GoogleDriveGetFileNode{
		Node: nodes.Node{
			ID:          nodeID,
			Version:     "0.0.1",
			Name:        "Google Drive Get File",
			Description: "Retrieve a file from Google Drive along with its metadata.",
			Icon: nodes.NodeIcon{
				Brand: "google_drive",
				Glyph: "file",
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
				"get",
				"download",
				"storage",
			},
			SupportedCredentials: []string{
				"google_drive_oauth2",
			},
			InputSchema: &gjs.Schema{
				Type: "object",
				Properties: map[string]*gjs.Schema{
					"fileId": {
						Type:        "string",
						Title:       "File ID",
						Description: "Identifier of the Google Drive file to retrieve.",
					},
				},
				Required: []string{
					"fileId",
				},
			},
			OutputSchema: &gjs.Schema{
				Type: "array",
				Items: &gjs.Schema{
					Type: "object",
					Properties: map[string]*gjs.Schema{
						"file": {
							Type:        "string",
							Description: "File contents (text) or hosted URL for binary uploads.",
						},
						"metadata": {
							Type: "object",
							Properties: map[string]*gjs.Schema{
								"id": {
									Type:        "string",
									Description: "Identifier of the file.",
								},
								"name": {
									Type:        "string",
									Description: "Name of the file.",
								},
								"size": {
									Type:        "string",
									Description: "Size of the file in bytes.",
								},
								"mimeType": {
									Type:        "string",
									Description: "MIME type of the file.",
								},
							},
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
