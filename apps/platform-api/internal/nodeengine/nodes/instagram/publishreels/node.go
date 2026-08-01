package publishreels

import (
	"github.com/blocknextai/go-packages/json"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
	gjs "github.com/google/jsonschema-go/jsonschema"
)

type InstagramPublishReelsNode struct {
	nodes.Node
}

func NewInstagramPublishReelsNode(nodeID string) *InstagramPublishReelsNode {
	return &InstagramPublishReelsNode{
		Node: nodes.Node{
			ID:          nodeID,
			Version:     "0.0.1",
			Name:        "Instagram Publish Reels",
			Description: "Publish a Reel to an Instagram account.",
			Icon: nodes.NodeIcon{
				Light: nodeID,
				Dark:  nodeID,
			},
			Categories:    []string{"Publishing"},
			SubCategories: []string{"Instagram"},
			Tags: []string{
				"social",
				"social-media",
				"video",
				"short-video",
				"publish",
				"share",
			},
			SupportedCredentials: []string{
				"instagram_oauth2",
			},
			InputSchema: &gjs.Schema{
				Type: "object",
				Properties: map[string]*gjs.Schema{
					"videoUrl": {
						Type:        "string",
						Title:       "Video URL",
						Description: "URL of the video to publish as a reel.",
						Format:      "uri",
					},
					"caption": {
						Type:        "string",
						Title:       "Caption",
						Description: "Caption text for the reel.",
					},
					"coverImageUrl": {
						Type:        "string",
						Title:       "Cover Image URL",
						Description: "Optional cover image URL for the reel.",
						Format:      "uri",
					},
					"shareToFeed": {
						Type:        "boolean",
						Title:       "Share to Feed",
						Description: "Whether to also share the reel to the Instagram feed.",
						Default:     json.RawMessage(`true`),
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
							Description: "Whether the reel was published successfully.",
						},
						"reelsId": {
							Type:        "string",
							Description: "Identifier of the published reel.",
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
