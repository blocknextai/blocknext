package createmusic

import (
	"github.com/blocknextai/go-packages/json"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
	gjs "github.com/google/jsonschema-go/jsonschema"
)

type SunoMusicCreateMusicNode struct {
	nodes.Node
}

func NewSunoMusicCreateMusicNode(nodeID string) *SunoMusicCreateMusicNode {
	return &SunoMusicCreateMusicNode{
		ID:          nodeID,
		Kind:        nodes.NodeKindAction,
		Version:     "0.0.1",
		Name:        "Suno Create Music",
		Description: "Generate a music track using Suno.",
		Icon: nodes.NodeIcon{
			Brand: "sunomusic",
			Glyph: "music",
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
			"audio",
			"generation",
			"sound",
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
					Description: "Prompt describing the music to generate.",
					MaxLength:   new(400),
				},
				"style": {
					Type:        "string",
					Title:       "Style",
					Description: "Musical style or genre for the generated track.",
					Default:     json.RawMessage(`"pop"`),
				},
				"title": {
					Type:        "string",
					Title:       "Title",
					Description: "Title for the generated music track.",
					Default:     json.RawMessage(`"via BlockNext"`),
				},
				"customMode": {
					Type:        "boolean",
					Title:       "Custom Mode",
					Description: "Whether to enable Suno custom mode.",
					Default:     json.RawMessage(`false`),
				},
				"instrumental": {
					Type:        "boolean",
					Title:       "Instrumental",
					Description: "Whether the generated track should be instrumental only.",
					Default:     json.RawMessage(`false`),
				},
				"negativeTags": {
					Type:        "string",
					Title:       "Negative Tags",
					Description: "Tags to exclude from the generated track.",
					Default:     json.RawMessage(`"blacklisted"`),
				},
				"model": {
					Type:        "string",
					Title:       "Model",
					Description: "Suno model version used for music generation.",
					Enum: []any{
						"V5_5",
						"V5",
						"V4_5ALL",
						"V4_5PLUS",
						"V4_5",
						"V4",
					},
					Default: json.RawMessage(`"V5"`),
				},
			},
			Required: []string{
				"prompt",
				"customMode",
				"instrumental",
				"negativeTags",
				"model",
			},
		},
		OutputSchema: &gjs.Schema{
			Type: "array",
			Items: &gjs.Schema{
				Type: "object",
				Properties: map[string]*gjs.Schema{
					"audio": {
						Type:        "string",
						Description: "URL of the generated audio track.",
						Format:      "uri",
					},
					"image": {
						Type:        "string",
						Description: "URL of the generated cover image.",
						Format:      "uri",
					},
				},
			},
		},
		HasNaturalLanguage: true,
		Annotations: nodes.NodeAnnotations{
			ReadOnly:    true,
			Destructive: new(false),
		},
	}
}
