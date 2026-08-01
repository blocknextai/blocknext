package publishstory

import (
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
	gjs "github.com/google/jsonschema-go/jsonschema"
)

type FacebookPublishStoryNode struct {
	nodes.Node
}

func NewFacebookPublishStoryNode(nodeID string) *FacebookPublishStoryNode {
	return &FacebookPublishStoryNode{
		Node: nodes.Node{
			ID:          nodeID,
			Version:     "0.0.1",
			Name:        "Facebook Publish Story",
			Description: "Publish a story to a Facebook page.",
			Icon: nodes.NodeIcon{
				Light: nodeID,
				Dark:  nodeID,
			},
			Categories:    []string{"Publishing"},
			SubCategories: []string{"Facebook"},
			Tags: []string{
				"social",
				"social-media",
				"publish",
				"story",
				"share",
				"media",
			},
			SupportedCredentials: []string{
				"facebook_oauth2",
			},
			InputSchema: &gjs.Schema{
				Type: "object",
				Properties: map[string]*gjs.Schema{
					"mediaUrls": {
						Type:        "array",
						Title:       "Media URLs",
						Description: "List of image URLs to publish as story frames.",
						Items: &gjs.Schema{
							Type:   "string",
							Format: "uri",
						},
					},
					"accountName": {
						Type:        "string",
						Title:       "Account Name",
						Description: "Name of the Facebook page to publish from.",
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
							Description: "Identifier of the published Facebook story.",
						},
						"storyUrl": {
							Type:        "string",
							Description: "Public URL of the published Facebook story.",
							Format:      "uri",
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
		},
	}
}
