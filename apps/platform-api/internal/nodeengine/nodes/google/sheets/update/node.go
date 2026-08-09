package update

import (
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
	gjs "github.com/google/jsonschema-go/jsonschema"
)

type GoogleSheetsUpdateNode struct {
	nodes.Node
}

func NewGoogleSheetsUpdateNode(nodeID string) *GoogleSheetsUpdateNode {
	return &GoogleSheetsUpdateNode{
		Node: nodes.Node{
			ID:          nodeID,
			Kind:        nodes.NodeKindAction,
			Version:     "0.0.1",
			Name:        "Google Sheets Update Cell",
			Description: "Update a single cell in a Google Sheets spreadsheet by row number and column.",
			Icon: nodes.NodeIcon{
				Brand: "google_sheets",
				Glyph: "pencil",
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
				"update",
				"modify",
				"data",
			},
			SupportedCredentials: []string{
				"google_sheets_oauth2",
			},
			InputSchema: &gjs.Schema{
				Type: "object",
				Properties: map[string]*gjs.Schema{
					"sheetId": {
						Type:        "string",
						Title:       "Sheet ID",
						Description: "Identifier of the Google Sheets spreadsheet to update.",
					},
					"rowNumber": {
						Type:        "number",
						Title:       "Row Number",
						Description: "1-based row number of the cell to update.",
					},
					"column": {
						Type:        "string",
						Title:       "Column",
						Description: "Column header name or column letter (e.g. \"name\" or \"B\").",
					},
					"newValue": {
						Type:        "string",
						Title:       "New Value",
						Description: "New cell value to write.",
					},
				},
				Required: []string{
					"sheetId",
					"rowNumber",
					"column",
					"newValue",
				},
			},
			OutputSchema: &gjs.Schema{
				Type: "array",
				Items: &gjs.Schema{
					Type: "object",
					Properties: map[string]*gjs.Schema{
						"status": {
							Type:        "boolean",
							Description: "Whether the cell was updated successfully.",
						},
						"sheetId": {
							Type:        "string",
							Description: "Identifier of the updated spreadsheet.",
						},
						"updatedRange": {
							Type:        "string",
							Description: "Range that was updated, in A1 notation.",
						},
						"rowNumber": {
							Type:        "integer",
							Description: "Row number that was updated.",
						},
						"column": {
							Type:        "string",
							Description: "Column that was updated.",
						},
						"newValue": {
							Type:        "string",
							Description: "Value that was written to the cell.",
						},
						"updatedCells": {
							Type:        "integer",
							Description: "Number of cells updated.",
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
