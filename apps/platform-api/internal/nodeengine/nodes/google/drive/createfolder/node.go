package createfolder

import (
	"github.com/blocknextai/go-packages/json"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
	gjs "github.com/google/jsonschema-go/jsonschema"
)

type GoogleDriveCreateFolderNode struct {
	nodes.Node
}

func NewGoogleDriveCreateFolderNode(nodeID string) *GoogleDriveCreateFolderNode {
	return &GoogleDriveCreateFolderNode{
		ID:          nodeID,
		Kind:        nodes.NodeKindAction,
		Version:     "0.0.1",
		Name:        "Google Drive Create Folder",
		Description: "Create a new folder in Google Drive.",
		Icon: nodes.NodeIcon{
			Brand: "google_drive",
			Glyph: "folder",
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
			"folder",
			"create",
			"storage",
			"file",
		},
		SupportedCredentials: []string{
			"google_drive_oauth2",
		},
		InputSchema: &gjs.Schema{
			Type: "object",
			Properties: map[string]*gjs.Schema{
				"folderName": {
					Type:        "string",
					Title:       "Folder Name",
					Description: "Name of the folder to create.",
				},
				"parentFolderId": {
					Type:        "string",
					Title:       "Parent Folder ID",
					Description: "Identifier of the parent folder; use \"root\" for the drive root.",
					Default:     json.RawMessage(`"root"`),
				},
			},
			Required: []string{
				"folderName",
			},
		},
		OutputSchema: &gjs.Schema{
			Type: "array",
			Items: &gjs.Schema{
				Type: "object",
				Properties: map[string]*gjs.Schema{
					"status": {
						Type:        "boolean",
						Description: "Whether the folder was created successfully.",
					},
					"folderId": {
						Type:        "string",
						Description: "Identifier of the created folder.",
					},
					"folderName": {
						Type:        "string",
						Description: "Name of the created folder.",
					},
					"webViewLink": {
						Type:        "string",
						Format:      "uri",
						Description: "Web view link to open the folder in Google Drive.",
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
