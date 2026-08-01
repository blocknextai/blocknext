package createrecord

import (
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
	gjs "github.com/google/jsonschema-go/jsonschema"
)

type AirtableCreateRecordNode struct {
	nodes.Node
}

func NewAirtableCreateRecordNode(nodeID string) *AirtableCreateRecordNode {
	return &AirtableCreateRecordNode{
		Node: nodes.Node{
			ID:          nodeID,
			Version:     "0.0.1",
			Name:        "Airtable Create Record",
			Description: "Create a new record in an Airtable base.",
			Icon: nodes.NodeIcon{
				Light: nodeID,
				Dark:  nodeID,
			},
			Categories:    []string{"Database"},
			SubCategories: []string{"Airtable"},
			Tags: []string{
				"database",
				"create",
				"record",
				"data",
				"productivity",
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
					"fields": {
						Type:        "string",
						Title:       "Fields",
						Description: "JSON string containing the fields and values for the new record.",
					},
				},
				Required: []string{
					"baseId",
					"tableId",
					"fields",
				},
			},
			OutputSchema: &gjs.Schema{
				Type: "array",
				Items: &gjs.Schema{
					Type: "object",
					Properties: map[string]*gjs.Schema{
						"status": {
							Type:        "boolean",
							Description: "Whether the record was created successfully.",
						},
						"id": {
							Type:        "string",
							Description: "Identifier of the newly created record.",
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
