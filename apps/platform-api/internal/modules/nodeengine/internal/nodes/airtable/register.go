package airtable

import (
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/application/jsonschema"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/mcp"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/providers"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/nodes/airtable/createrecord"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/nodes/airtable/listrecords"
)

func Register() {
	nodeID := "airtable"

	createRecordNodeID := nodeID + "_create_record"
	createRecordNode := createrecord.NewAirtableCreateRecordNode(createRecordNodeID)
	createRecordValidator := jsonschema.New[createrecord.AirtableCreateRecordExecutorInput](createRecordNode.GetInputSchema())
	createRecordExecutor := createrecord.NewAirtableCreateRecordExecutor(createRecordNodeID, createRecordValidator)

	listRecordsNodeID := nodeID + "_list_records"
	listRecordsNode := listrecords.NewAirtableListRecordsNode(listRecordsNodeID)
	listRecordsValidator := jsonschema.New[listrecords.AirtableListRecordsExecutorInput](listRecordsNode.GetInputSchema())
	listRecordsExecutor := listrecords.NewAirtableListRecordsExecutor(listRecordsNodeID, listRecordsValidator)

	providers.Register(providers.Provider{
		Server: mcp.Server{
			ID:          nodeID,
			Name:        "Airtable",
			Description: "Tools for managing Airtable records.",
			Icon: mcp.ServerIcon{
				Brand: "airtable",
			},
			Version: "0.0.1",
		},
		Actions: []providers.Action{
			{Node: createRecordNode, Executor: createRecordExecutor},
			{Node: listRecordsNode, Executor: listRecordsExecutor},
		},
	})
}
