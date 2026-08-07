package publishmediapost

import (
	"github.com/blocknextai/go-packages/json"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
	gjs "github.com/google/jsonschema-go/jsonschema"
)

type LinkedinPublishMediaPostNode struct {
	nodes.Node
}

func NewLinkedinPublishMediaPostNode(nodeID string) *LinkedinPublishMediaPostNode {
	return &LinkedinPublishMediaPostNode{
		Node: nodes.Node{
			ID:          nodeID,
			Version:     "0.0.1",
			Name:        "LinkedIn Publish Media Post",
			Description: "Publish a post with media to a LinkedIn account.",
			Icon: nodes.NodeIcon{
				Brand: "linkedin",
				Glyph: "image",
			},
			Inputs: []nodes.NodeHandle{
				{Key: "in"},
			},
			Outputs: []nodes.NodeHandle{
				{Key: "out"},
			},
			Categories:    []string{"Publishing"},
			SubCategories: []string{"LinkedIn"},
			Tags: []string{
				"social",
				"social-media",
				"professional",
				"b2b",
				"publish",
				"media",
				"share",
			},
			SupportedCredentials: []string{
				"linkedin_oauth2",
			},
			InputSchema: &gjs.Schema{
				Type: "object",
				Properties: map[string]*gjs.Schema{
					"text": {
						Type:        "string",
						Title:       "Text",
						Description: "Text content of the LinkedIn post.",
					},
					"mediaUrl": {
						Type:        "string",
						Title:       "Media URL",
						Description: "URL of the media (image or video) to attach to the post.",
						Format:      "uri",
					},
					"mediaType": {
						Type:        "string",
						Title:       "Media Type",
						Description: "Type of media being attached (IMAGE or VIDEO).",
						Enum:        []any{"IMAGE", "VIDEO"},
						Default:     json.RawMessage(`"IMAGE"`),
					},
					"visibility": {
						Type:        "string",
						Title:       "Visibility",
						Description: "Visibility of the post (PUBLIC or CONNECTIONS).",
						Enum:        []any{"PUBLIC", "CONNECTIONS"},
						Default:     json.RawMessage(`"PUBLIC"`),
					},
				},
				Required: []string{
					"mediaUrl",
					"mediaType",
					"visibility",
				},
			},
			OutputSchema: &gjs.Schema{
				Type: "array",
				Items: &gjs.Schema{
					Type: "object",
					Properties: map[string]*gjs.Schema{
						"status": {
							Type:        "boolean",
							Description: "Whether the media post was published successfully.",
						},
						"postId": {
							Type:        "string",
							Description: "Identifier of the published LinkedIn post.",
						},
						"postUrl": {
							Type:        "string",
							Description: "URL of the published LinkedIn post.",
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
