package sendmessage

import (
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
	gjs "github.com/google/jsonschema-go/jsonschema"
)

type SlackSendMessageNode struct {
	nodes.Node
}

func NewSlackSendMessageNode(nodeID string) *SlackSendMessageNode {
	return &SlackSendMessageNode{
		Node: nodes.Node{
			ID:          nodeID,
			Version:     "0.0.1",
			Name:        "Slack Send Message",
			Description: "Send a text message to a Slack channel.",
			Icon: nodes.NodeIcon{
				Brand: "slack",
				Glyph: "chat",
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
				"message",
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
						Description: "Slack channel ID where the message will be posted.",
					},
					"text": {
						Type:        "string",
						Title:       "Text",
						Description: "Message text to post to the Slack channel.",
					},
				},
				Required: []string{
					"channel",
					"text",
				},
			},
			OutputSchema: &gjs.Schema{
				Type: "array",
				Items: &gjs.Schema{
					Type: "object",
					Properties: map[string]*gjs.Schema{
						"status": {
							Type:        "boolean",
							Description: "Whether the Slack message was sent successfully.",
						},
						"messageId": {
							Type:        "string",
							Description: "Slack message timestamp identifier.",
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
