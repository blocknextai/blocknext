package drive

import (
	"github.com/blocknextai/platform-api/internal/filegateway"
	"github.com/blocknextai/platform-api/internal/nodeengine/application/jsonschema"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/executors"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/functioncalling"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/mcp"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/google/drive/createfolder"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/google/drive/createtextfile"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/google/drive/getfile"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/google/drive/uploadfile"
)

func Register(fileGateway filegateway.FileGateway) {
	nodeID := "google_drive"

	createFolderNodeID := nodeID + "_create_folder"
	createFolderNode := createfolder.NewGoogleDriveCreateFolderNode(createFolderNodeID)
	createFolderValidator := jsonschema.New[createfolder.GoogleDriveCreateFolderExecutorInput](createFolderNode.GetInputSchema())
	createFolderExecutor := createfolder.NewGoogleDriveCreateFolderExecutor(createFolderNodeID, createFolderValidator)

	nodes.RegisterNode(createFolderNode)
	executors.RegisterExecutor(createFolderExecutor)
	functioncalling.RegisterFunctionCalling(functioncalling.Generate(createFolderNode))

	createTextFileNodeID := nodeID + "_create_text_file"
	createTextFileNode := createtextfile.NewGoogleDriveCreateTextFileNode(createTextFileNodeID)
	createTextFileValidator := jsonschema.New[createtextfile.GoogleDriveCreateTextFileExecutorInput](createTextFileNode.GetInputSchema())
	createTextFileExecutor := createtextfile.NewGoogleDriveCreateTextFileExecutor(createTextFileNodeID, createTextFileValidator)

	nodes.RegisterNode(createTextFileNode)
	executors.RegisterExecutor(createTextFileExecutor)
	functioncalling.RegisterFunctionCalling(functioncalling.Generate(createTextFileNode))

	getFileNodeID := nodeID + "_get_file"
	getFileNode := getfile.NewGoogleDriveGetFileNode(getFileNodeID)
	getFileValidator := jsonschema.New[getfile.GoogleDriveGetFileExecutorInput](getFileNode.GetInputSchema())
	getFileExecutor := getfile.NewGoogleDriveGetFileExecutor(getFileNodeID, getFileValidator, fileGateway)

	nodes.RegisterNode(getFileNode)
	executors.RegisterExecutor(getFileExecutor)
	functioncalling.RegisterFunctionCalling(functioncalling.Generate(getFileNode))

	uploadFileNodeID := nodeID + "_upload_file"
	uploadFileNode := uploadfile.NewGoogleDriveUploadFileNode(uploadFileNodeID)
	uploadFileValidator := jsonschema.New[uploadfile.GoogleDriveUploadFileExecutorInput](uploadFileNode.GetInputSchema())
	uploadFileExecutor := uploadfile.NewGoogleDriveUploadFileExecutor(uploadFileNodeID, uploadFileValidator, fileGateway)

	nodes.RegisterNode(uploadFileNode)
	executors.RegisterExecutor(uploadFileExecutor)
	functioncalling.RegisterFunctionCalling(functioncalling.Generate(uploadFileNode))

	mcp.RegisterServer(&mcp.Server{
		ID:          nodeID,
		Name:        "Google Drive",
		Description: "Tools for managing files and folders in Google Drive.",
		Version:     "0.0.1",
		Tools: []nodes.NodeManager{
			createFolderNode,
			createTextFileNode,
			getFileNode,
			uploadFileNode,
		},
	})
}
