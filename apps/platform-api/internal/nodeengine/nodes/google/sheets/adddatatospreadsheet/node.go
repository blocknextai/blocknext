package adddatatospreadsheet

import (
	"github.com/blocknextai/go-packages/json"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
	gjs "github.com/google/jsonschema-go/jsonschema"
)

type GoogleSheetsAddDataToSpreadsheetNode struct {
	nodes.Node
}

func NewGoogleSheetsAddDataToSpreadsheetNode(nodeID string) *GoogleSheetsAddDataToSpreadsheetNode {
	return &GoogleSheetsAddDataToSpreadsheetNode{
		ID:          nodeID,
		Kind:        nodes.NodeKindAction,
		Name:        "Google Sheets Add Data to Spreadsheet",
		Version:     "0.0.1",
		Description: "Append rows of data to a Google Sheets spreadsheet.",
		Icon: nodes.NodeIcon{
			Brand: "google_sheets",
			Glyph: "plus",
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
			"add",
			"append",
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
					Description: "Google Sheets spreadsheet identifier.",
				},
				"range": {
					Type:        "string",
					Title:       "Range",
					Description: "A1 notation range to append data to (e.g. Sheet1!A1).",
				},
				"data": {
					Type:        "array",
					Title:       "Data",
					Description: "Two-dimensional array of rows and cell values to append.",
					Items: &gjs.Schema{
						Type:  "array",
						Items: &gjs.Schema{Type: "string"},
					},
				},
				"valueInputOption": {
					Type:        "string",
					Title:       "Value Input Option",
					Description: "How input data should be interpreted by Google Sheets.",
					Enum:        []any{"USER_ENTERED", "RAW"},
					Default:     json.RawMessage(`"USER_ENTERED"`),
				},
			},
			Required: []string{
				"spreadsheetId",
				"range",
				"data",
			},
		},
		OutputSchema: &gjs.Schema{
			Type: "array",
			Items: &gjs.Schema{
				Type: "object",
				Properties: map[string]*gjs.Schema{
					"status": {
						Type:        "boolean",
						Description: "Whether the data was appended successfully.",
					},
					"spreadsheetId": {
						Type:        "string",
						Description: "Identifier of the updated spreadsheet.",
					},
					"updatedRange": {
						Type:        "string",
						Description: "Range that was updated, in A1 notation.",
					},
					"updatedRows": {
						Type:        "integer",
						Description: "Number of rows updated.",
					},
					"updatedColumns": {
						Type:        "integer",
						Description: "Number of columns updated.",
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
			Destructive: new(false),
		},
	}
}
