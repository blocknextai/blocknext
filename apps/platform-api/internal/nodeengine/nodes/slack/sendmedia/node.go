package sendmedia

import (
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
	gjs "github.com/google/jsonschema-go/jsonschema"
)

type SlackSendMediaNode struct {
	nodes.Node
}

func NewSlackSendMediaNode(nodeID string) *SlackSendMediaNode {
	return &SlackSendMediaNode{
		Node: nodes.Node{
			ID:          nodeID,
			Kind:        nodes.NodeKindAction,
			Version:     "0.0.1",
			Name:        "Slack Send Media",
			Description: "Send a media file to a Slack channel.",
			Icon: nodes.NodeIcon{
				Brand: "slack",
				Glyph: "image",
			},
			Inputs: []nodes.NodeHandle{
				{Key: "in"},
			},
			Outputs: []nodes.NodeHandle{
				{Key: "out"},
			},
			Categories:    []string{"Publishing"},
			SubCategories: []string{"Slack"},
			Tags: []string{
				"chat",
				"messaging",
				"send",
				"media",
				"team",
				"workspace",
				"notification",
			},
			SupportedCredentials: []string{
				"slack_oauth2",
			},
			InputSchema: &gjs.Schema{
				Type: "object",
				Properties: map[string]*gjs.Schema{
					"channel": {
						Type:        "string",
						Title:       "Channel",
						Description: "Slack channel ID where the media will be sent.",
					},
					"mediaUrls": {
						Type:        "array",
						Title:       "Media URLs",
						Description: "List of media URLs to upload to the Slack channel.",
						Items:       &gjs.Schema{Type: "string"},
					},
					"text": {
						Type:        "string",
						Title:       "Text",
						Description: "Optional comment to post alongside the media.",
					},
				},
				Required: []string{
					"channel",
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
							Description: "Whether the media upload completed successfully.",
						},
						"error": {
							Type:        "string",
							Description: "Error message when the media upload failed.",
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
