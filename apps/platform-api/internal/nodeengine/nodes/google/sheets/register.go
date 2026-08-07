package sheets

import (
	"github.com/blocknextai/platform-api/internal/nodeengine/application/jsonschema"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/executors"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/functioncalling"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/mcp"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/google/sheets/adddatatospreadsheet"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/google/sheets/createspreadsheet"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/google/sheets/deletespreadsheet"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/google/sheets/readspreadsheet"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/google/sheets/update"
)

func Register() {
	nodeID := "sheets"

	createSpreadsheetNodeID := nodeID + "_create_spreadsheet"
	createSpreadsheetNode := createspreadsheet.NewGoogleSheetsCreateSpreadsheetNode(createSpreadsheetNodeID)
	createSpreadsheetValidator := jsonschema.New[createspreadsheet.GoogleSheetsCreateSpreadsheetExecutorInput](createSpreadsheetNode.GetInputSchema())
	createSpreadsheetExecutor := createspreadsheet.NewGoogleSheetsCreateSpreadsheetExecutor(createSpreadsheetNodeID, createSpreadsheetValidator)

	nodes.RegisterNode(createSpreadsheetNode)
	executors.RegisterExecutor(createSpreadsheetExecutor)
	functioncalling.RegisterFunctionCalling(functioncalling.Generate(createSpreadsheetNode))

	readSpreadsheetNodeID := nodeID + "_read_spreadsheet"
	readSpreadsheetNode := readspreadsheet.NewGoogleSheetsReadSpreadsheetNode(readSpreadsheetNodeID)
	readSpreadsheetValidator := jsonschema.New[readspreadsheet.GoogleSheetsReadSpreadsheetExecutorInput](readSpreadsheetNode.GetInputSchema())
	readSpreadsheetExecutor := readspreadsheet.NewGoogleSheetsReadSpreadsheetExecutor(readSpreadsheetNodeID, readSpreadsheetValidator)

	nodes.RegisterNode(readSpreadsheetNode)
	executors.RegisterExecutor(readSpreadsheetExecutor)
	functioncalling.RegisterFunctionCalling(functioncalling.Generate(readSpreadsheetNode))

	deleteSpreadsheetNodeID := nodeID + "_delete_spreadsheet"
	deleteSpreadsheetNode := deletespreadsheet.NewGoogleSheetsDeleteSpreadsheetNode(deleteSpreadsheetNodeID)
	deleteSpreadsheetValidator := jsonschema.New[deletespreadsheet.GoogleSheetsDeleteSpreadsheetExecutorInput](deleteSpreadsheetNode.GetInputSchema())
	deleteSpreadsheetExecutor := deletespreadsheet.NewGoogleSheetsDeleteSpreadsheetExecutor(deleteSpreadsheetNodeID, deleteSpreadsheetValidator)

	nodes.RegisterNode(deleteSpreadsheetNode)
	executors.RegisterExecutor(deleteSpreadsheetExecutor)
	functioncalling.RegisterFunctionCalling(functioncalling.Generate(deleteSpreadsheetNode))

	addDataToSpreadsheetNodeID := nodeID + "_add_data_to_spreadsheet"
	addDataToSpreadsheetNode := adddatatospreadsheet.NewGoogleSheetsAddDataToSpreadsheetNode(addDataToSpreadsheetNodeID)
	addDataToSpreadsheetValidator := jsonschema.New[adddatatospreadsheet.GoogleSheetsAddDataToSpreadsheetExecutorInput](addDataToSpreadsheetNode.GetInputSchema())
	addDataToSpreadsheetExecutor := adddatatospreadsheet.NewGoogleSheetsAddDataToSpreadsheetExecutor(addDataToSpreadsheetNodeID, addDataToSpreadsheetValidator)

	nodes.RegisterNode(addDataToSpreadsheetNode)
	executors.RegisterExecutor(addDataToSpreadsheetExecutor)
	functioncalling.RegisterFunctionCalling(functioncalling.Generate(addDataToSpreadsheetNode))

	updateNodeID := nodeID + "_update"
	updateNode := update.NewGoogleSheetsUpdateNode(updateNodeID)
	updateValidator := jsonschema.New[update.GoogleSheetsUpdateExecutorInput](updateNode.GetInputSchema())
	updateExecutor := update.NewGoogleSheetsUpdateExecutor(updateNodeID, updateValidator)

	nodes.RegisterNode(updateNode)
	executors.RegisterExecutor(updateExecutor)
	functioncalling.RegisterFunctionCalling(functioncalling.Generate(updateNode))

	mcp.RegisterServer(&mcp.Server{
		ID:          nodeID,
		Name:        "Google Sheets",
		Description: "Tools for creating, reading, updating, and deleting Google Sheets data.",
		Icon: mcp.ServerIcon{
			Brand: "google_sheets",
		},
		Version: "0.0.1",
		Tools: []nodes.NodeManager{
			createSpreadsheetNode,
			readSpreadsheetNode,
			deleteSpreadsheetNode,
			addDataToSpreadsheetNode,
			updateNode,
		},
	})
}
