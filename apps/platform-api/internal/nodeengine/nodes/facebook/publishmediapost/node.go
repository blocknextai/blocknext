package publishmediapost

import (
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
	gjs "github.com/google/jsonschema-go/jsonschema"
)

type FacebookPublishMediaPostNode struct {
	nodes.Node
}

func NewFacebookPublishMediaPostNode(nodeID string) *FacebookPublishMediaPostNode {
	return &FacebookPublishMediaPostNode{
		Node: nodes.Node{
			ID:          nodeID,
			Kind:        nodes.NodeKindAction,
			Version:     "0.0.1",
			Name:        "Facebook Publish Media Post",
			Description: "Publish a post with media to a Facebook page.",
			Icon: nodes.NodeIcon{
				Brand: "facebook",
				Glyph: "image",
			},
			Inputs: []nodes.NodeHandle{
				{Key: "in"},
			},
			Outputs: []nodes.NodeHandle{
				{Key: "out"},
			},
			Categories:    []string{"Publishing"},
			SubCategories: []string{"Facebook"},
			Tags: []string{
				"social",
				"social-media",
				"publish",
				"media",
				"post",
				"share",
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
						Description: "List of image or video URLs to attach to the post.",
						Items: &gjs.Schema{
							Type:   "string",
							Format: "uri",
						},
					},
					"message": {
						Type:        "string",
						Title:       "Message",
						Description: "Caption or message text for the post.",
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
							Description: "Whether the post was published successfully.",
						},
						"postId": {
							Type:        "string",
							Description: "Identifier of the published Facebook post.",
						},
						"postUrl": {
							Type:        "string",
							Description: "Public URL of the published Facebook post.",
							Format:      "uri",
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
