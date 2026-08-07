package keep

import (
	"github.com/blocknextai/platform-api/internal/nodeengine/application/jsonschema"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/executors"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/functioncalling"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/mcp"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/google/keep/createnote"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/google/keep/deletenote"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/google/keep/getnote"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/google/keep/listnotes"
)

func Register() {
	nodeID := "google_keep"

	createNoteNodeID := nodeID + "_create_note"
	createNoteNode := createnote.NewGoogleKeepCreateNoteNode(createNoteNodeID)
	createNoteValidator := jsonschema.New[createnote.GoogleKeepCreateNoteExecutorInput](createNoteNode.GetInputSchema())
	createNoteExecutor := createnote.NewGoogleKeepCreateNoteExecutor(createNoteNodeID, createNoteValidator)

	nodes.RegisterNode(createNoteNode)
	executors.RegisterExecutor(createNoteExecutor)
	functioncalling.RegisterFunctionCalling(functioncalling.Generate(createNoteNode))

	getNoteNodeID := nodeID + "_get_note"
	getNoteNode := getnote.NewGoogleKeepGetNoteNode(getNoteNodeID)
	getNoteValidator := jsonschema.New[getnote.GoogleKeepGetNoteExecutorInput](getNoteNode.GetInputSchema())
	getNoteExecutor := getnote.NewGoogleKeepGetNoteExecutor(getNoteNodeID, getNoteValidator)

	nodes.RegisterNode(getNoteNode)
	executors.RegisterExecutor(getNoteExecutor)
	functioncalling.RegisterFunctionCalling(functioncalling.Generate(getNoteNode))

	listNotesNodeID := nodeID + "_list_notes"
	listNotesNode := listnotes.NewGoogleKeepListNotesNode(listNotesNodeID)
	listNotesValidator := jsonschema.New[listnotes.GoogleKeepListNotesExecutorInput](listNotesNode.GetInputSchema())
	listNotesExecutor := listnotes.NewGoogleKeepListNotesExecutor(listNotesNodeID, listNotesValidator)

	nodes.RegisterNode(listNotesNode)
	executors.RegisterExecutor(listNotesExecutor)
	functioncalling.RegisterFunctionCalling(functioncalling.Generate(listNotesNode))

	deleteNoteNodeID := nodeID + "_delete_note"
	deleteNoteNode := deletenote.NewGoogleKeepDeleteNoteNode(deleteNoteNodeID)
	deleteNoteValidator := jsonschema.New[deletenote.GoogleKeepDeleteNoteExecutorInput](deleteNoteNode.GetInputSchema())
	deleteNoteExecutor := deletenote.NewGoogleKeepDeleteNoteExecutor(deleteNoteNodeID, deleteNoteValidator)

	nodes.RegisterNode(deleteNoteNode)
	executors.RegisterExecutor(deleteNoteExecutor)
	functioncalling.RegisterFunctionCalling(functioncalling.Generate(deleteNoteNode))

	mcp.RegisterServer(&mcp.Server{
		ID:          nodeID,
		Name:        "Google Keep",
		Description: "Tools for creating, reading, listing, and deleting Google Keep notes.",
		Icon: mcp.ServerIcon{
			Brand: "google_keep",
		},
		Version: "0.0.1",
		Tools: []nodes.NodeManager{
			createNoteNode,
			getNoteNode,
			listNotesNode,
			deleteNoteNode,
		},
	})
}
