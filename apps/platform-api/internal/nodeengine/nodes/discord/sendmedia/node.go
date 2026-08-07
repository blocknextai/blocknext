package sendmedia

import (
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
	gjs "github.com/google/jsonschema-go/jsonschema"
)

type DiscordSendMediaNode struct {
	nodes.Node
}

func NewDiscordSendMediaNode(nodeID string) *DiscordSendMediaNode {
	return &DiscordSendMediaNode{
		Node: nodes.Node{
			ID:          nodeID,
			Version:     "0.0.1",
			Name:        "Discord Send Media",
			Description: "Send a media file to a Discord channel.",
			Icon: nodes.NodeIcon{
				Brand: "discord",
				Glyph: "image",
			},
			Inputs: []nodes.NodeHandle{
				{Key: "in"},
			},
			Outputs: []nodes.NodeHandle{
				{Key: "out"},
			},
			Categories:    []string{"Publishing"},
			SubCategories: []string{"Discord"},
			Tags: []string{
				"chat",
				"messaging",
				"send",
				"media",
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
						Description: "Discord channel identifier where the media will be sent.",
					},
					"mediaUrls": {
						Type:        "array",
						Title:       "Media URLs",
						Description: "List of media file URLs to upload to the channel.",
						Items: &gjs.Schema{
							Type:   "string",
							Format: "uri",
						},
					},
				},
				Required: []string{
					"channelId",
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
							Description: "Whether the media was sent successfully.",
						},
						"messageId": {
							Type:        "string",
							Description: "Identifier of the Discord message that was created.",
						},
						"error": {
							Type:        "string",
							Description: "Error message if sending failed.",
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
