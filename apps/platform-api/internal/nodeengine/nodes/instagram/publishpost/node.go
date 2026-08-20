package publishpost

import (
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
	gjs "github.com/google/jsonschema-go/jsonschema"
)

type InstagramPublishPostNode struct {
	nodes.Node
}

func NewInstagramPublishPostNode(nodeID string) *InstagramPublishPostNode {
	return &InstagramPublishPostNode{
		ID:          nodeID,
		Kind:        nodes.NodeKindAction,
		Version:     "0.0.1",
		Name:        "Instagram Publish Post",
		Description: "Publish a post to an Instagram account.",
		Icon: nodes.NodeIcon{
			Brand: "instagram",
			Glyph: "send",
		},
		Inputs: []nodes.NodeHandle{
			{Key: "in"},
		},
		Outputs: []nodes.NodeHandle{
			{Key: "out"},
		},
		Categories:    []string{"Publishing"},
		SubCategories: []string{"Instagram"},
		Tags: []string{
			"social",
			"social-media",
			"publish",
			"post",
			"share",
			"media",
		},
		SupportedCredentials: []string{
			"instagram_oauth2",
		},
		InputSchema: &gjs.Schema{
			Type: "object",
			Properties: map[string]*gjs.Schema{
				"caption": {
					Type:        "string",
					Title:       "Caption",
					Description: "Caption text for the Instagram post.",
				},
				"mediaUrls": {
					Type:        "array",
					Title:       "Media URLs",
					Description: "List of image or video URLs to publish.",
					Items: &gjs.Schema{
						Type:   "string",
						Format: "uri",
					},
				},
				"location": {
					Type:        "string",
					Title:       "Location",
					Description: "Optional Instagram location identifier.",
				},
				"altText": {
					Type:        "string",
					Title:       "Alt Text",
					Description: "Optional accessibility alt text for the media.",
				},
			},
			Required: []string{
				"mediaUrls",
			},
		},
		OutputSchema: &gjs.Schema{
			Type: "array",
			Items: &gjs.Schema{
				Type: "object",
				Properties: map[string]*gjs.Schema{
					"status": {
						Type:        "boolean",
						Description: "Whether the post was published successfully.",
					},
					"postId": {
						Type:        "string",
						Description: "Identifier of the published Instagram post.",
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
