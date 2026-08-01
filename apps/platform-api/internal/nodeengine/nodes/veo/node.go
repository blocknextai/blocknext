package veo

import (
	"github.com/blocknextai/go-packages/json"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
	gjs "github.com/google/jsonschema-go/jsonschema"
)

type VeoNode struct {
	nodes.Node
}

func NewVeoNode(nodeID string) *VeoNode {
	return &VeoNode{
		Node: nodes.Node{
			ID:          nodeID,
			Version:     "0.0.1",
			Name:        "Veo",
			Description: "Generate a video from a text prompt using Veo.",
			Icon: nodes.NodeIcon{
				Light: nodeID,
				Dark:  nodeID,
			},
			Categories:    []string{"Video"},
			SubCategories: []string{"Veo"},
			Tags: []string{
				"google",
				"ai",
				"video",
				"generation",
				"media",
				"creative",
			},
			SupportedCredentials: []string{
				"veo_api",
			},
			InputSchema: &gjs.Schema{
				Type: "object",
				Properties: map[string]*gjs.Schema{
					"model": {
						Type:        "string",
						Title:       "Model",
						Description: "Veo model identifier to use for generation.",
						Enum: []any{
							"veo-3.1-generate-preview",
							"veo-3.1-lite-generate-preview",
						},
						Default: json.RawMessage(`"veo-3.1-lite-generate-preview"`),
					},
					"prompt": {
						Type:        "string",
						Title:       "Prompt",
						Description: "Text prompt describing the video to generate.",
					},
					"image": {
						Type:        "string",
						Title:       "Image",
						Description: "Optional reference image URL for image-to-video generation.",
					},
					"aspectRatio": {
						Type:        "string",
						Title:       "Aspect Ratio",
						Description: "Aspect ratio of the generated video.",
						Enum:        []any{"16:9"},
						Default:     json.RawMessage(`"16:9"`),
					},
				},
			},
			OutputSchema: &gjs.Schema{
				Type: "array",
				Items: &gjs.Schema{
					Type: "object",
					Properties: map[string]*gjs.Schema{
						"video": {
							Type:        "string",
							Description: "URL of the generated Veo video.",
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
