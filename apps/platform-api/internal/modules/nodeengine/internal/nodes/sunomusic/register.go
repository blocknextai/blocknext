package sunomusic

import (
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/application/jsonschema"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/mcp"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/providers"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/nodes/sunomusic/createlyrics"
	"github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/nodes/sunomusic/createmusic"
)

func Register() {
	nodeID := "sunomusic"

	createMusicNodeID := nodeID + "_create_music"
	createMusicNode := createmusic.NewSunoMusicCreateMusicNode(createMusicNodeID)
	createMusicValidator := jsonschema.New[createmusic.SunoMusicCreateMusicInput](createMusicNode.GetInputSchema())
	createMusicExecutor := createmusic.NewSunoMusicCreateMusicExecutor(createMusicNodeID, createMusicValidator)

	createLyricsNodeID := nodeID + "_create_lyrics"
	createLyricsNode := createlyrics.NewSunoMusicCreateLyricsNode(createLyricsNodeID)
	createLyricsValidator := jsonschema.New[createlyrics.SunoMusicCreateLyricsInput](createLyricsNode.GetInputSchema())
	createLyricsExecutor := createlyrics.NewSunoMusicCreateLyricsExecutor(createLyricsNodeID, createLyricsValidator)

	providers.Register(providers.Provider{
		Server: mcp.Server{
			ID:          nodeID,
			Name:        "Suno Music",
			Description: "Tools for generating music and lyrics via Suno.",
			Icon: mcp.ServerIcon{
				Brand: "sunomusic",
			},
			Version: "0.0.1",
		},
		Actions: []providers.Action{
			{Node: createMusicNode, Executor: createMusicExecutor},
			{Node: createLyricsNode, Executor: createLyricsExecutor},
		},
	})
}
