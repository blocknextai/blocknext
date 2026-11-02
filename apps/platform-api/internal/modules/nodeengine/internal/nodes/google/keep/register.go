package keep

import (
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/application/jsonschema"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/mcp"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/providers"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/nodes/google/keep/createnote"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/nodes/google/keep/deletenote"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/nodes/google/keep/getnote"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/nodes/google/keep/listnotes"
)

func Register() {
	nodeID := "google_keep"

	createNoteNodeID := nodeID + "_create_note"
	createNoteNode := createnote.NewGoogleKeepCreateNoteNode(createNoteNodeID)
	createNoteValidator := jsonschema.New[createnote.GoogleKeepCreateNoteExecutorInput](createNoteNode.GetInputSchema())
	createNoteExecutor := createnote.NewGoogleKeepCreateNoteExecutor(createNoteNodeID, createNoteValidator)

	getNoteNodeID := nodeID + "_get_note"
	getNoteNode := getnote.NewGoogleKeepGetNoteNode(getNoteNodeID)
	getNoteValidator := jsonschema.New[getnote.GoogleKeepGetNoteExecutorInput](getNoteNode.GetInputSchema())
	getNoteExecutor := getnote.NewGoogleKeepGetNoteExecutor(getNoteNodeID, getNoteValidator)

	listNotesNodeID := nodeID + "_list_notes"
	listNotesNode := listnotes.NewGoogleKeepListNotesNode(listNotesNodeID)
	listNotesValidator := jsonschema.New[listnotes.GoogleKeepListNotesExecutorInput](listNotesNode.GetInputSchema())
	listNotesExecutor := listnotes.NewGoogleKeepListNotesExecutor(listNotesNodeID, listNotesValidator)

	deleteNoteNodeID := nodeID + "_delete_note"
	deleteNoteNode := deletenote.NewGoogleKeepDeleteNoteNode(deleteNoteNodeID)
	deleteNoteValidator := jsonschema.New[deletenote.GoogleKeepDeleteNoteExecutorInput](deleteNoteNode.GetInputSchema())
	deleteNoteExecutor := deletenote.NewGoogleKeepDeleteNoteExecutor(deleteNoteNodeID, deleteNoteValidator)

	providers.Register(providers.Provider{
		Server: mcp.Server{
			ID:          nodeID,
			Name:        "Google Keep",
			Description: "Tools for creating, reading, listing, and deleting Google Keep notes.",
			Icon: mcp.ServerIcon{
				Brand: "google_keep",
			},
			Version: "0.0.1",
		},
		Actions: []providers.Action{
			{Node: createNoteNode, Executor: createNoteExecutor},
			{Node: getNoteNode, Executor: getNoteExecutor},
			{Node: listNotesNode, Executor: listNotesExecutor},
			{Node: deleteNoteNode, Executor: deleteNoteExecutor},
		},
	})
}
