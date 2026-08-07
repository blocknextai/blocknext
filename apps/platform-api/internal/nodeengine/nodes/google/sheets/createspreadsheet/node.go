package createspreadsheet

import (
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
	gjs "github.com/google/jsonschema-go/jsonschema"
)

type GoogleSheetsCreateSpreadsheetNode struct {
	nodes.Node
}

func NewGoogleSheetsCreateSpreadsheetNode(nodeID string) *GoogleSheetsCreateSpreadsheetNode {
	return &GoogleSheetsCreateSpreadsheetNode{
		Node: nodes.Node{
			ID:          nodeID,
			Version:     "0.0.1",
			Name:        "Google Sheets Create Spreadsheet",
			Description: "Create a new Google Sheets spreadsheet with the given title.",
			Icon: nodes.NodeIcon{
				Brand: "google_sheets",
				Glyph: "table",
			},
			Inputs: []nodes.NodeHandle{
				{Key: "in"},
			},
			Outputs: []nodes.NodeHandle{
				{Key: "out"},
			},
			Categories:    []string{"Google Workspace"},
			SubCategories: []string{"Google Sheets"},
			Tags: []string{
				"google",
				"sheets",
				"spreadsheet",
				"create",
				"data",
				"productivity",
			},
			SupportedCredentials: []string{
				"google_sheets_oauth2",
			},
			InputSchema: &gjs.Schema{
				Type: "object",
				Properties: map[string]*gjs.Schema{
					"title": {
						Type:        "string",
						Title:       "Title",
						Description: "Title of the new spreadsheet.",
					},
				},
				Required: []string{
					"title",
				},
			},
			OutputSchema: &gjs.Schema{
				Type: "array",
				Items: &gjs.Schema{
					Type: "object",
					Properties: map[string]*gjs.Schema{
						"status": {
							Type:        "boolean",
							Description: "Whether the spreadsheet was created successfully.",
						},
						"id": {
							Type:        "string",
							Description: "Identifier of the created spreadsheet.",
						},
						"name": {
							Type:        "string",
							Description: "Name of the created spreadsheet.",
						},
						"webViewLink": {
							Type:        "string",
							Format:      "uri",
							Description: "Web view link to open the spreadsheet.",
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
