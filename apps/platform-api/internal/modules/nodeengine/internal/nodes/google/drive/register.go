package drive

import (
	"github.com/blocknextai/platform-api/internal/filegateway"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/application/jsonschema"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/mcp"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/providers"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/nodes/google/drive/createfolder"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/nodes/google/drive/createtextfile"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/nodes/google/drive/getfile"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/nodes/google/drive/uploadfile"
)

func Register(fileGateway filegateway.FileGateway) {
	nodeID := "google_drive"

	createFolderNodeID := nodeID + "_create_folder"
	createFolderNode := createfolder.NewGoogleDriveCreateFolderNode(createFolderNodeID)
	createFolderValidator := jsonschema.New[createfolder.GoogleDriveCreateFolderExecutorInput](createFolderNode.GetInputSchema())
	createFolderExecutor := createfolder.NewGoogleDriveCreateFolderExecutor(createFolderNodeID, createFolderValidator)

	createTextFileNodeID := nodeID + "_create_text_file"
	createTextFileNode := createtextfile.NewGoogleDriveCreateTextFileNode(createTextFileNodeID)
	createTextFileValidator := jsonschema.New[createtextfile.GoogleDriveCreateTextFileExecutorInput](createTextFileNode.GetInputSchema())
	createTextFileExecutor := createtextfile.NewGoogleDriveCreateTextFileExecutor(createTextFileNodeID, createTextFileValidator)

	getFileNodeID := nodeID + "_get_file"
	getFileNode := getfile.NewGoogleDriveGetFileNode(getFileNodeID)
	getFileValidator := jsonschema.New[getfile.GoogleDriveGetFileExecutorInput](getFileNode.GetInputSchema())
	getFileExecutor := getfile.NewGoogleDriveGetFileExecutor(getFileNodeID, getFileValidator, fileGateway)

	uploadFileNodeID := nodeID + "_upload_file"
	uploadFileNode := uploadfile.NewGoogleDriveUploadFileNode(uploadFileNodeID)
	uploadFileValidator := jsonschema.New[uploadfile.GoogleDriveUploadFileExecutorInput](uploadFileNode.GetInputSchema())
	uploadFileExecutor := uploadfile.NewGoogleDriveUploadFileExecutor(uploadFileNodeID, uploadFileValidator, fileGateway)

	providers.Register(providers.Provider{
		Server: mcp.Server{
			ID:          nodeID,
			Name:        "Google Drive",
			Description: "Tools for managing files and folders in Google Drive.",
			Icon: mcp.ServerIcon{
				Brand: "google_drive",
			},
			Version: "0.0.1",
		},
		Actions: []providers.Action{
			{Node: createFolderNode, Executor: createFolderExecutor},
			{Node: createTextFileNode, Executor: createTextFileExecutor},
			{Node: getFileNode, Executor: getFileExecutor},
			{Node: uploadFileNode, Executor: uploadFileExecutor},
		},
	})
}
