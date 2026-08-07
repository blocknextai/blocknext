package createlyrics

import (
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
	gjs "github.com/google/jsonschema-go/jsonschema"
)

type SunoMusicCreateLyricsNode struct {
	nodes.Node
}

func NewSunoMusicCreateLyricsNode(nodeID string) *SunoMusicCreateLyricsNode {
	return &SunoMusicCreateLyricsNode{
		Node: nodes.Node{
			ID:          nodeID,
			Version:     "0.0.1",
			Name:        "Suno Create Lyrics",
			Description: "Generate song lyrics using Suno.",
			Icon: nodes.NodeIcon{
				Brand: "sunomusic",
				Glyph: "note",
			},
			Inputs: []nodes.NodeHandle{
				{Key: "in"},
			},
			Outputs: []nodes.NodeHandle{
				{Key: "out"},
			},
			Categories:    []string{"Audio"},
			SubCategories: []string{"Suno"},
			Tags: []string{
				"ai",
				"music",
				"lyrics",
				"text",
				"generation",
				"creative",
			},
			SupportedCredentials: []string{
				"sunomusic_api",
			},
			InputSchema: &gjs.Schema{
				Type: "object",
				Properties: map[string]*gjs.Schema{
					"prompt": {
						Type:        "string",
						Title:       "Prompt",
						Description: "Prompt describing the lyrics to generate.",
					},
				},
				Required: []string{
					"prompt",
				},
			},
			OutputSchema: &gjs.Schema{
				Type: "array",
				Items: &gjs.Schema{
					Type: "object",
					Properties: map[string]*gjs.Schema{
						"text": {
							Type:        "string",
							Description: "Generated lyrics text.",
						},
					},
				},
			},
			HasNaturalLanguage: true,
			Annotations: nodes.NodeAnnotations{
				ReadOnly:    true,
				Destructive: new(false),
			},
		},
	}
}
