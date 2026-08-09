package readspreadsheet

import (
	"github.com/blocknextai/go-packages/json"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
	gjs "github.com/google/jsonschema-go/jsonschema"
)

type GoogleSheetsReadSpreadsheetNode struct {
	nodes.Node
}

func NewGoogleSheetsReadSpreadsheetNode(nodeID string) *GoogleSheetsReadSpreadsheetNode {
	return &GoogleSheetsReadSpreadsheetNode{
		Node: nodes.Node{
			ID:          nodeID,
			Kind:        nodes.NodeKindAction,
			Version:     "0.0.1",
			Name:        "Google Sheets Read Spreadsheet",
			Description: "Read rows from a Google Sheets spreadsheet, optionally filtering and limiting results.",
			Icon: nodes.NodeIcon{
				Brand: "google_sheets",
				Glyph: "eye",
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
				"read",
				"get",
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
						Description: "Identifier of the Google Sheets spreadsheet to read.",
					},
					"range": {
						Type:        "string",
						Title:       "Range",
						Description: "A1 notation range to read (e.g. A:Z).",
						Default:     json.RawMessage(`"A:Z"`),
					},
					"filterColumn": {
						Type:        "string",
						Title:       "Filter Column",
						Description: "Column name used to filter returned rows.",
					},
					"filterValue": {
						Type:        "string",
						Title:       "Filter Value",
						Description: "Value the filter column must match.",
					},
					"limit": {
						Type:        "number",
						Title:       "Limit",
						Description: "Maximum number of rows to return; 0 means no limit.",
					},
					"headerRow": {
						Type:        "number",
						Title:       "Header Row",
						Description: "Row number containing column headers; 0 means no header row.",
						Default:     json.RawMessage(`0`),
					},
				},
				Required: []string{
					"spreadsheetId",
				},
			},
			OutputSchema: &gjs.Schema{
				Type: "array",
				Items: &gjs.Schema{
					Type:        "object",
					Description: "One object per spreadsheet row. Property names come from the header row (or column_N when no header is set), plus row_number.",
					Properties: map[string]*gjs.Schema{
						"row_number": {
							Type:        "integer",
							Description: "Row index in the source spreadsheet.",
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
