package nanobanana

import (
	"github.com/blocknextai/go-packages/json"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
	gjs "github.com/google/jsonschema-go/jsonschema"
)

type GeminiImageGenerationNode struct {
	nodes.Node
}

func NewGeminiImageGenerationNode(nodeID string) *GeminiImageGenerationNode {
	return &GeminiImageGenerationNode{
		Node: nodes.Node{
			ID:          nodeID,
			Kind:        nodes.NodeKindAction,
			Version:     "0.0.1",
			Name:        "Gemini Nano Banana",
			Description: "Generate an image from a text prompt using Google Gemini Nano Banana.",
			Icon: nodes.NodeIcon{
				Brand: "gemini",
				Glyph: "image",
			},
			Inputs: []nodes.NodeHandle{
				{Key: "in"},
			},
			Outputs: []nodes.NodeHandle{
				{Key: "out"},
			},
			Categories:    []string{"Image"},
			SubCategories: []string{"Gemini"},
			Tags: []string{
				"image_generation",
				"google",
				"ai",
				"image",
				"generation",
				"media",
				"creative",
			},
			SupportedCredentials: []string{
				"gemini_api",
			},
			InputSchema: &gjs.Schema{
				Type: "object",
				Properties: map[string]*gjs.Schema{
					"prompt": {
						Type:        "string",
						Title:       "Prompt",
						Description: "Text prompt describing the image to generate.",
					},
					"model": {
						Type:        "string",
						Title:       "Model",
						Description: "Gemini image-generation model identifier.",
						Enum: []any{
							"gemini-3-pro-image",
							"gemini-3.1-flash-image",
							"gemini-3.1-flash-lite-image",
						},
						Default: json.RawMessage(`"gemini-3.1-flash-image"`),
					},
					"image": {
						Type:        "string",
						Title:       "Image",
						Description: "Optional reference image URL to guide generation.",
						Format:      "uri",
					},
				},
				Required: []string{
					"prompt",
					"model",
				},
			},
			OutputSchema: &gjs.Schema{
				Type: "array",
				Items: &gjs.Schema{
					Type: "object",
					Properties: map[string]*gjs.Schema{
						"image": {
							Type:        "string",
							Description: "URL of the generated image.",
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
		},
	}
}
