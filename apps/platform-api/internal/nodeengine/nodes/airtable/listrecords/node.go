package listrecords

import (
	"github.com/blocknextai/go-packages/json"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
	gjs "github.com/google/jsonschema-go/jsonschema"
)

type AirtableListRecordsNode struct {
	nodes.Node
}

func NewAirtableListRecordsNode(nodeID string) *AirtableListRecordsNode {
	return &AirtableListRecordsNode{
		Node: nodes.Node{
			ID:          nodeID,
			Kind:        nodes.NodeKindAction,
			Version:     "0.0.1",
			Name:        "Airtable List Records",
			Description: "List records from an Airtable base.",
			Icon: nodes.NodeIcon{
				Brand: "airtable",
				Glyph: "list",
			},
			Inputs: []nodes.NodeHandle{
				{Key: "in"},
			},
			Outputs: []nodes.NodeHandle{
				{Key: "out"},
			},
			Categories:    []string{"Database"},
			SubCategories: []string{"Airtable"},
			Tags: []string{
				"database",
				"list",
				"records",
				"data",
				"query",
				"spreadsheet",
			},
			SupportedCredentials: []string{
				"airtable_api",
			},
			InputSchema: &gjs.Schema{
				Type: "object",
				Properties: map[string]*gjs.Schema{
					"baseId": {
						Type:        "string",
						Title:       "Base ID",
						Description: "Identifier of the Airtable base.",
					},
					"tableId": {
						Type:        "string",
						Title:       "Table ID",
						Description: "Identifier or name of the Airtable table.",
					},
					"maxRecords": {
						Type:        "number",
						Title:       "Max Records",
						Description: "Maximum number of records to return.",
						Maximum:     new(100.0),
						Default:     json.RawMessage(`100`),
					},
				},
				Required: []string{
					"baseId",
					"tableId",
				},
			},
			OutputSchema: &gjs.Schema{
				Type: "array",
				Items: &gjs.Schema{
					Type: "object",
					Properties: map[string]*gjs.Schema{
						"id": {
							Type:        "string",
							Description: "Identifier of the Airtable record.",
						},
						"createdTime": {
							Type:        "string",
							Description: "Timestamp of when the record was created.",
						},
						"fields": {
							Type:        "object",
							Description: "Field values of the Airtable record.",
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
