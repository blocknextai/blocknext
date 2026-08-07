package publishpost

import (
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
	gjs "github.com/google/jsonschema-go/jsonschema"
)

type FacebookPublishPostNode struct {
	nodes.Node
}

func NewFacebookPublishPostNode(nodeID string) *FacebookPublishPostNode {
	return &FacebookPublishPostNode{
		Node: nodes.Node{
			ID:          nodeID,
			Version:     "0.0.1",
			Name:        "Facebook Publish Post",
			Description: "Publish a text post to a Facebook page.",
			Icon: nodes.NodeIcon{
				Brand: "facebook",
				Glyph: "send",
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
				"post",
				"share",
			},
			SupportedCredentials: []string{
				"facebook_oauth2",
			},
			InputSchema: &gjs.Schema{
				Type: "object",
				Properties: map[string]*gjs.Schema{
					"message": {
						Type:        "string",
						Title:       "Message",
						Description: "Text content of the Facebook post.",
					},
					"accountName": {
						Type:        "string",
						Title:       "Account Name",
						Description: "Name of the Facebook page to publish from.",
					},
				},
				Required: []string{
					"message",
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
