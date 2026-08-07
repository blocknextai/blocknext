package sunomusic

import (
	"github.com/blocknextai/platform-api/internal/nodeengine/application/jsonschema"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/executors"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/functioncalling"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/mcp"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/sunomusic/createlyrics"
	"github.com/blocknextai/platform-api/internal/nodeengine/nodes/sunomusic/createmusic"
)

func Register() {
	nodeID := "sunomusic"

	createMusicNodeID := nodeID + "_create_music"
	createMusicNode := createmusic.NewSunoMusicCreateMusicNode(createMusicNodeID)
	createMusicValidator := jsonschema.New[createmusic.SunoMusicCreateMusicInput](createMusicNode.GetInputSchema())
	createMusicExecutor := createmusic.NewSunoMusicCreateMusicExecutor(createMusicNodeID, createMusicValidator)

	nodes.RegisterNode(createMusicNode)
	executors.RegisterExecutor(createMusicExecutor)
	functioncalling.RegisterFunctionCalling(functioncalling.Generate(createMusicNode))

	createLyricsNodeID := nodeID + "_create_lyrics"
	createLyricsNode := createlyrics.NewSunoMusicCreateLyricsNode(createLyricsNodeID)
	createLyricsValidator := jsonschema.New[createlyrics.SunoMusicCreateLyricsInput](createLyricsNode.GetInputSchema())
	createLyricsExecutor := createlyrics.NewSunoMusicCreateLyricsExecutor(createLyricsNodeID, createLyricsValidator)

	nodes.RegisterNode(createLyricsNode)
	executors.RegisterExecutor(createLyricsExecutor)
	functioncalling.RegisterFunctionCalling(functioncalling.Generate(createLyricsNode))

	mcp.RegisterServer(&mcp.Server{
		ID:          nodeID,
		Name:        "Suno Music",
		Description: "Tools for generating music and lyrics via Suno.",
		Icon: mcp.ServerIcon{
			Brand: "sunomusic",
		},
		Version: "0.0.1",
		Tools: []nodes.NodeManager{
			createMusicNode,
			createLyricsNode,
		},
	})
}
