package publishpost

import (
	"github.com/blocknextai/go-packages/json"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
	gjs "github.com/google/jsonschema-go/jsonschema"
)

type XPublishPostNode struct {
	nodes.Node
}

func NewXPublishPostNode(nodeID string) *XPublishPostNode {
	return &XPublishPostNode{
		ID:          nodeID,
		Kind:        nodes.NodeKindAction,
		Version:     "0.0.1",
		Name:        "X Publish Post",
		Description: "Publish a post to an X (Twitter) account.",
		Icon: nodes.NodeIcon{
			Brand: "x",
			Glyph: "send",
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
			"post",
			"tweet",
			"share",
		},
		SupportedCredentials: []string{
			"x_oauth2",
		},
		InputSchema: &gjs.Schema{
			Type: "object",
			Properties: map[string]*gjs.Schema{
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
				"forSuperFollowersOnly": {
					Type:        "boolean",
					Title:       "For Super Followers Only",
					Description: "Whether the tweet is visible to super followers only.",
					Default:     json.RawMessage(`false`),
				},
				"possiblyContentSensitive": {
					Type:        "boolean",
					Title:       "Possibly Content Sensitive",
					Description: "Whether the tweet may contain sensitive content.",
					Default:     json.RawMessage(`false`),
				},
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
