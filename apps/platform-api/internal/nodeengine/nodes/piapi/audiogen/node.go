package audiogen

import (
	"github.com/blocknextai/go-packages/json"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
	gjs "github.com/google/jsonschema-go/jsonschema"
)

type PiAPIAudioGenNode struct {
	nodes.Node
}

func NewPiAPIAudioGenNode(nodeID string) *PiAPIAudioGenNode {
	return &PiAPIAudioGenNode{
		Node: nodes.Node{
			ID:          nodeID,
			Version:     "0.0.1",
			Name:        "PiAPI Audio Generation",
			Description: "Generate audio using PiAPI.",
			Icon: nodes.NodeIcon{
				Brand: "piapi",
				Glyph: "speaker",
			},
			Inputs: []nodes.NodeHandle{
				{Key: "in"},
			},
			Outputs: []nodes.NodeHandle{
				{Key: "out"},
			},
			Categories:    []string{"Audio"},
			SubCategories: []string{"PiApi"},
			Tags: []string{
				"ai",
				"audio",
				"music",
				"generation",
				"sound",
				"creative",
			},
			SupportedCredentials: []string{
				"piapi_api",
			},
			InputSchema: &gjs.Schema{
				Type: "object",
				Properties: map[string]*gjs.Schema{
					"prompt": {
						Type:        "string",
						Title:       "Prompt",
						Description: "Prompt describing the audio to generate.",
					},
					"negativeTags": {
						Type:        "string",
						Title:       "Negative Tags",
						Description: "Tags to exclude from the generated audio.",
					},
					"gptDescriptionPrompt": {
						Type:        "string",
						Title:       "GPT Description Prompt",
						Description: "Optional GPT-generated description prompt for guidance.",
					},
					"title": {
						Type:        "string",
						Title:       "Title",
						Description: "Title for the generated track.",
					},
					"lyricsType": {
						Type:        "string",
						Title:       "Lyrics Type",
						Description: "How lyrics should be produced for the track.",
						Enum:        []any{"generate", "user", "instrumental"},
						Default:     json.RawMessage(`"generate"`),
					},
					"seed": {
						Type:        "number",
						Title:       "Seed",
						Description: "Seed used for deterministic generation.",
					},
					"lyrics": {
						Type:        "string",
						Title:       "Lyrics",
						Description: "User-provided lyrics when lyricsType is set to user.",
					},
				},
				Required: []string{
					"prompt",
					"lyricsType",
				},
			},
			OutputSchema: &gjs.Schema{
				Type: "array",
				Items: &gjs.Schema{
					Type: "object",
					Properties: map[string]*gjs.Schema{
						"audio": {
							Type:        "string",
							Description: "URL of the generated audio file.",
							Format:      "uri",
						},
						"image": {
							Type:        "string",
							Description: "URL of the generated cover image.",
							Format:      "uri",
						},
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
