package uploadfile

import (
	"github.com/blocknextai/go-packages/json"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
	gjs "github.com/google/jsonschema-go/jsonschema"
)

type GoogleDriveUploadFileNode struct {
	nodes.Node
}

func NewGoogleDriveUploadFileNode(nodeID string) *GoogleDriveUploadFileNode {
	return &GoogleDriveUploadFileNode{
		ID:          nodeID,
		Kind:        nodes.NodeKindAction,
		Version:     "0.0.1",
		Name:        "Google Drive Upload File",
		Description: "Upload one or more files to Google Drive from a list of source URLs.",
		Icon: nodes.NodeIcon{
			Brand: "google_drive",
			Glyph: "upload",
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
			"upload",
			"storage",
			"share",
		},
		SupportedCredentials: []string{
			"google_drive_oauth2",
		},
		InputSchema: &gjs.Schema{
			Type: "object",
			Properties: map[string]*gjs.Schema{
				"folderId": {
					Type:        "string",
					Title:       "Folder ID",
					Description: "Identifier of the destination folder; use \"root\" for the drive root.",
					Default:     json.RawMessage(`"root"`),
				},
				"files": {
					Type:        "array",
					Title:       "Files",
					Description: "List of file URLs to upload to Google Drive.",
					Items: &gjs.Schema{
						Type:   "string",
						Format: "uri",
					},
				},
			},
			Required: []string{
				"files",
			},
		},
		OutputSchema: &gjs.Schema{
			Type: "array",
			Items: &gjs.Schema{
				Type: "object",
				Properties: map[string]*gjs.Schema{
					"status": {
						Type:        "boolean",
						Description: "Whether all files were uploaded successfully.",
					},
				},
			},
		},
		HasNaturalLanguage: true,
		Annotations: nodes.NodeAnnotations{
			Destructive: new(false),
		},
	}
}
