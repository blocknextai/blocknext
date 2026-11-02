package sheets

import (
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/application/jsonschema"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/mcp"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/providers"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/nodes/google/sheets/adddatatospreadsheet"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/nodes/google/sheets/createspreadsheet"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/nodes/google/sheets/deletespreadsheet"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/nodes/google/sheets/readspreadsheet"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/nodes/google/sheets/update"
)

func Register() {
	nodeID := "sheets"

	createSpreadsheetNodeID := nodeID + "_create_spreadsheet"
	createSpreadsheetNode := createspreadsheet.NewGoogleSheetsCreateSpreadsheetNode(createSpreadsheetNodeID)
	createSpreadsheetValidator := jsonschema.New[createspreadsheet.GoogleSheetsCreateSpreadsheetExecutorInput](createSpreadsheetNode.GetInputSchema())
	createSpreadsheetExecutor := createspreadsheet.NewGoogleSheetsCreateSpreadsheetExecutor(createSpreadsheetNodeID, createSpreadsheetValidator)

	readSpreadsheetNodeID := nodeID + "_read_spreadsheet"
	readSpreadsheetNode := readspreadsheet.NewGoogleSheetsReadSpreadsheetNode(readSpreadsheetNodeID)
	readSpreadsheetValidator := jsonschema.New[readspreadsheet.GoogleSheetsReadSpreadsheetExecutorInput](readSpreadsheetNode.GetInputSchema())
	readSpreadsheetExecutor := readspreadsheet.NewGoogleSheetsReadSpreadsheetExecutor(readSpreadsheetNodeID, readSpreadsheetValidator)

	deleteSpreadsheetNodeID := nodeID + "_delete_spreadsheet"
	deleteSpreadsheetNode := deletespreadsheet.NewGoogleSheetsDeleteSpreadsheetNode(deleteSpreadsheetNodeID)
	deleteSpreadsheetValidator := jsonschema.New[deletespreadsheet.GoogleSheetsDeleteSpreadsheetExecutorInput](deleteSpreadsheetNode.GetInputSchema())
	deleteSpreadsheetExecutor := deletespreadsheet.NewGoogleSheetsDeleteSpreadsheetExecutor(deleteSpreadsheetNodeID, deleteSpreadsheetValidator)

	addDataToSpreadsheetNodeID := nodeID + "_add_data_to_spreadsheet"
	addDataToSpreadsheetNode := adddatatospreadsheet.NewGoogleSheetsAddDataToSpreadsheetNode(addDataToSpreadsheetNodeID)
	addDataToSpreadsheetValidator := jsonschema.New[adddatatospreadsheet.GoogleSheetsAddDataToSpreadsheetExecutorInput](addDataToSpreadsheetNode.GetInputSchema())
	addDataToSpreadsheetExecutor := adddatatospreadsheet.NewGoogleSheetsAddDataToSpreadsheetExecutor(addDataToSpreadsheetNodeID, addDataToSpreadsheetValidator)

	updateNodeID := nodeID + "_update"
	updateNode := update.NewGoogleSheetsUpdateNode(updateNodeID)
	updateValidator := jsonschema.New[update.GoogleSheetsUpdateExecutorInput](updateNode.GetInputSchema())
	updateExecutor := update.NewGoogleSheetsUpdateExecutor(updateNodeID, updateValidator)

	providers.Register(providers.Provider{
		Server: mcp.Server{
			ID:          nodeID,
			Name:        "Google Sheets",
			Description: "Tools for creating, reading, updating, and deleting Google Sheets data.",
			Icon: mcp.ServerIcon{
				Brand: "google_sheets",
			},
			Version: "0.0.1",
		},
		Actions: []providers.Action{
			{Node: createSpreadsheetNode, Executor: createSpreadsheetExecutor},
			{Node: readSpreadsheetNode, Executor: readSpreadsheetExecutor},
			{Node: deleteSpreadsheetNode, Executor: deleteSpreadsheetExecutor},
			{Node: addDataToSpreadsheetNode, Executor: addDataToSpreadsheetExecutor},
			{Node: updateNode, Executor: updateExecutor},
		},
	})
}
