package publishpost

import (
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
	gjs "github.com/google/jsonschema-go/jsonschema"
)

type TiktokPublishPostNode struct {
	nodes.Node
}

func NewTiktokPublishPostNode(nodeID string) *TiktokPublishPostNode {
	return &TiktokPublishPostNode{
		ID:          nodeID,
		Kind:        nodes.NodeKindAction,
		Version:     "0.0.1",
		Name:        "TikTok Publish Post",
		Description: "Publish a video post to a TikTok account.",
		Icon: nodes.NodeIcon{
			Brand: "tiktok",
			Glyph: "send",
		},
		Inputs: []nodes.NodeHandle{
			{Key: "in"},
		},
		Outputs: []nodes.NodeHandle{
			{Key: "out"},
		},
		Categories:    []string{"Publishing"},
		SubCategories: []string{"TikTok"},
		Tags: []string{
			"social",
			"social-media",
			"video",
			"short-video",
			"publish",
			"share",
		},
		SupportedCredentials: []string{
			"tiktok_oauth2",
		},
		InputSchema: &gjs.Schema{
			Type: "object",
			Properties: map[string]*gjs.Schema{
				"caption": {
					Type:        "string",
					Title:       "Caption",
					Description: "Caption text to attach to the TikTok post.",
				},
				"videoUrl": {
					Type:        "string",
					Title:       "Video URL",
					Description: "URL of the video to publish to TikTok.",
					Format:      "uri",
				},
			},
			Required: []string{
				"videoUrl",
			},
		},
		OutputSchema: &gjs.Schema{
			Type: "array",
			Items: &gjs.Schema{
				Type: "object",
				Properties: map[string]*gjs.Schema{
					"status": {
						Type:        "boolean",
						Description: "Whether the TikTok post was published successfully.",
					},
					"publishId": {
						Type:        "string",
						Description: "Identifier of the TikTok publish operation.",
					},
				},
			},
		},
		HasNaturalLanguage: true,
		Annotations: nodes.NodeAnnotations{
			Destructive: new(false),
		},
	}
}
