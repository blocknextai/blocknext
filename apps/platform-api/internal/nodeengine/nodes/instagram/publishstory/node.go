package publishstory

import (
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
	gjs "github.com/google/jsonschema-go/jsonschema"
)

type InstagramPublishStoryNode struct {
	nodes.Node
}

func NewInstagramPublishStoryNode(nodeID string) *InstagramPublishStoryNode {
	return &InstagramPublishStoryNode{
		ID:          nodeID,
		Kind:        nodes.NodeKindAction,
		Version:     "0.0.1",
		Name:        "Instagram Publish Story",
		Description: "Publish a story to an Instagram account.",
		Icon: nodes.NodeIcon{
			Brand: "instagram",
			Glyph: "story",
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
			"story",
			"share",
			"media",
		},
		SupportedCredentials: []string{
			"instagram_oauth2",
		},
		InputSchema: &gjs.Schema{
			Type: "object",
			Properties: map[string]*gjs.Schema{
				"mediaUrls": {
					Type:        "array",
					Title:       "Media URLs",
					Description: "List of image or video URLs to publish as stories.",
					Items: &gjs.Schema{
						Type:   "string",
						Format: "uri",
					},
				},
				"storyLink": {
					Type:        "string",
					Title:       "Story Link",
					Description: "Optional link to attach to the story.",
					Format:      "uri",
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
						Description: "Whether the story was published successfully.",
					},
					"storyId": {
						Type:        "string",
						Description: "Identifier of the published story.",
					},
					"error": {
						Type:        "string",
						Description: "Error message if publishing failed.",
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
