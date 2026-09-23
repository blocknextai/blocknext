package airtable

import (
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/application/jsonschema"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/executors"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/functioncalling"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/mcp"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/nodes"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/nodes/airtable/createrecord"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/nodes/airtable/listrecords"
)

func Register() {
	nodeID := "airtable"

	createRecordNodeID := nodeID + "_create_record"
	createRecordNode := createrecord.NewAirtableCreateRecordNode(createRecordNodeID)
	createRecordValidator := jsonschema.New[createrecord.AirtableCreateRecordExecutorInput](createRecordNode.GetInputSchema())
	createRecordExecutor := createrecord.NewAirtableCreateRecordExecutor(createRecordNodeID, createRecordValidator)

	nodes.RegisterNode(createRecordNode)
	executors.RegisterExecutor(createRecordExecutor)
	functioncalling.RegisterFunctionCalling(functioncalling.Generate(createRecordNode))

	listRecordsNodeID := nodeID + "_list_records"
	listRecordsNode := listrecords.NewAirtableListRecordsNode(listRecordsNodeID)
	listRecordsValidator := jsonschema.New[listrecords.AirtableListRecordsExecutorInput](listRecordsNode.GetInputSchema())
	listRecordsExecutor := listrecords.NewAirtableListRecordsExecutor(listRecordsNodeID, listRecordsValidator)

	nodes.RegisterNode(listRecordsNode)
	executors.RegisterExecutor(listRecordsExecutor)
	functioncalling.RegisterFunctionCalling(functioncalling.Generate(listRecordsNode))

	mcp.RegisterServer(&mcp.Server{
		ID:          nodeID,
		Name:        "Airtable",
		Description: "Tools for managing Airtable records.",
		Icon: mcp.ServerIcon{
			Brand: "airtable",
		},
		Version: "0.0.1",
		Tools: []nodes.NodeManager{
			createRecordNode,
			listRecordsNode,
		},
	})
}
