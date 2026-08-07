package soundcloud

import (
	"github.com/blocknextai/platform-api/internal/filegateway"
	"github.com/blocknextai/platform-api/internal/nodeengine/application/jsonschema"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/executors"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/functioncalling"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/mcp"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/soundcloud/createplaylist"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/soundcloud/createtrack"
)

func Register(fileGateway filegateway.FileGateway) {
	nodeID := "soundcloud"

	createTrackNodeID := nodeID + "_create_track"
	createTrackNode := createtrack.NewSoundCloudCreateTrackNode(createTrackNodeID)
	createTrackValidator := jsonschema.New[createtrack.SoundCloudCreateTrackExecutorInput](createTrackNode.GetInputSchema())
	createTrackExecutor := createtrack.NewSoundCloudCreateTrackExecutor(createTrackNodeID, createTrackValidator, fileGateway)

	nodes.RegisterNode(createTrackNode)
	executors.RegisterExecutor(createTrackExecutor)
	functioncalling.RegisterFunctionCalling(functioncalling.Generate(createTrackNode))

	createPlaylistNodeID := nodeID + "_create_playlist"
	createPlaylistNode := createplaylist.NewSoundCloudCreatePlaylistNode(createPlaylistNodeID)
	createPlaylistValidator := jsonschema.New[createplaylist.SoundCloudCreatePlaylistExecutorInput](createPlaylistNode.GetInputSchema())
	createPlaylistExecutor := createplaylist.NewSoundCloudCreatePlaylistExecutor(createPlaylistNodeID, createPlaylistValidator)

	nodes.RegisterNode(createPlaylistNode)
	executors.RegisterExecutor(createPlaylistExecutor)
	functioncalling.RegisterFunctionCalling(functioncalling.Generate(createPlaylistNode))

	mcp.RegisterServer(&mcp.Server{
		ID:          nodeID,
		Name:        "SoundCloud",
		Description: "Tools for managing SoundCloud tracks and playlists.",
		Icon: mcp.ServerIcon{
			Brand: "soundcloud",
		},
		Version: "0.0.1",
		Tools: []nodes.NodeManager{
			createTrackNode,
			createPlaylistNode,
		},
	})
}
