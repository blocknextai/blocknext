package publishmediapost

import (
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
	gjs "github.com/google/jsonschema-go/jsonschema"
)

type XPublishMediaPostNode struct {
	nodes.Node
}

func NewXPublishMediaPostNode(nodeID string) *XPublishMediaPostNode {
	return &XPublishMediaPostNode{
		ID:          nodeID,
		Kind:        nodes.NodeKindAction,
		Version:     "0.0.1",
		Name:        "X Publish Media Post",
		Description: "Publish a post with media to an X (Twitter) account.",
		Icon: nodes.NodeIcon{
			Brand: "x",
			Glyph: "image",
		},
		Inputs: []nodes.NodeHandle{
			{Key: "in"},
		},
		Outputs: []nodes.NodeHandle{
			{Key: "out"},
		},
		Categories:    []string{"Publishing"},
		SubCategories: []string{"X"},
		Tags: []string{
			"twitter",
			"social",
			"social-media",
			"publish",
			"media",
			"tweet",
			"share",
		},
		SupportedCredentials: []string{
			"x_oauth2",
		},
		InputSchema: &gjs.Schema{
			Type: "object",
			Properties: map[string]*gjs.Schema{
				"mediaUrls": {
					Type:        "array",
					Title:       "Media URLs",
					Description: "List of media URLs to attach to the tweet.",
					Items:       &gjs.Schema{Type: "string"},
				},
				"text": {
					Type:        "string",
					Title:       "Text",
					Description: "Text content of the tweet to publish.",
				},
				"replyToTweetId": {
					Type:        "string",
					Title:       "Reply To Tweet ID",
					Description: "Identifier of the tweet to reply to.",
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
						Description: "Whether the tweet was published successfully.",
					},
					"postId": {
						Type:        "string",
						Description: "Identifier of the published tweet.",
					},
					"postUrl": {
						Type:        "string",
						Description: "URL of the published tweet.",
						Format:      "uri",
					},
					"rateLimit": {
						Type: "object",
						Properties: map[string]*gjs.Schema{
							"appLimit24Hour":           {Type: "string"},
							"appLimit24HourRemaining":  {Type: "string"},
							"appLimit24HourReset":      {Type: "string"},
							"userLimit24Hour":          {Type: "string"},
							"userLimit24HourRemaining": {Type: "string"},
							"userLimit24HourReset":     {Type: "string"},
							"rateLimitLimit":           {Type: "string"},
							"rateLimitRemaining":       {Type: "string"},
							"rateLimitReset":           {Type: "string"},
						},
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
