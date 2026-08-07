package deletespreadsheet

import (
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
	gjs "github.com/google/jsonschema-go/jsonschema"
)

type GoogleSheetsDeleteSpreadsheetNode struct {
	nodes.Node
}

func NewGoogleSheetsDeleteSpreadsheetNode(nodeID string) *GoogleSheetsDeleteSpreadsheetNode {
	return &GoogleSheetsDeleteSpreadsheetNode{
		Node: nodes.Node{
			ID:          nodeID,
			Version:     "0.0.1",
			Name:        "Google Sheets Delete Spreadsheet",
			Description: "Delete a Google Sheets spreadsheet by its identifier.",
			Icon: nodes.NodeIcon{
				Brand: "google_sheets",
				Glyph: "trash",
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
				"delete",
				"remove",
				"data",
			},
			SupportedCredentials: []string{
				"google_sheets_oauth2",
			},
			InputSchema: &gjs.Schema{
				Type: "object",
				Properties: map[string]*gjs.Schema{
					"spreadsheetId": {
						Type:        "string",
						Title:       "Spreadsheet ID",
						Description: "Identifier of the Google Sheets spreadsheet to delete.",
					},
				},
				Required: []string{
					"spreadsheetId",
				},
			},
			OutputSchema: &gjs.Schema{
				Type: "array",
				Items: &gjs.Schema{
					Type: "object",
					Properties: map[string]*gjs.Schema{
						"status": {
							Type:        "boolean",
							Description: "Whether the spreadsheet was deleted successfully.",
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
