package sendmessage

import (
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
	gjs "github.com/google/jsonschema-go/jsonschema"
)

type DiscordSendMessageNode struct {
	nodes.Node
}

func NewDiscordSendMessageNode(nodeID string) *DiscordSendMessageNode {
	return &DiscordSendMessageNode{
		Node: nodes.Node{
			ID:          nodeID,
			Version:     "0.0.1",
			Name:        "Discord Send Message",
			Description: "Send a text message to a Discord channel.",
			Icon: nodes.NodeIcon{
				Light: nodeID,
				Dark:  nodeID,
			},
			Categories:    []string{"Publishing"},
			SubCategories: []string{"Discord"},
			Tags: []string{
				"chat",
				"messaging",
				"send",
				"message",
				"community",
				"notification",
			},
			SupportedCredentials: []string{
				"discord_api",
			},
			InputSchema: &gjs.Schema{
				Type: "object",
				Properties: map[string]*gjs.Schema{
					"channelId": {
						Type:        "string",
						Title:       "Channel ID",
						Description: "Discord channel identifier where the message will be sent.",
					},
					"content": {
						Type:        "string",
						Title:       "Content",
						Description: "Text content of the Discord message.",
					},
				},
				Required: []string{
					"channelId",
					"content",
				},
			},
			OutputSchema: &gjs.Schema{
				Type: "array",
				Items: &gjs.Schema{
					Type: "object",
					Properties: map[string]*gjs.Schema{
						"status": {
							Type:        "boolean",
							Description: "Whether the message was sent successfully.",
						},
						"messageId": {
							Type:        "string",
							Description: "Identifier of the Discord message that was created.",
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
