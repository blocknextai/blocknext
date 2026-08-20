package sendmedia

import (
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
	gjs "github.com/google/jsonschema-go/jsonschema"
)

type TelegramSendMediaNode struct {
	nodes.Node
}

func NewTelegramSendMediaNode(nodeID string) *TelegramSendMediaNode {
	return &TelegramSendMediaNode{
		ID:          nodeID,
		Kind:        nodes.NodeKindAction,
		Version:     "0.0.1",
		Name:        "Telegram Send Media",
		Description: "Send a media file to a Telegram chat.",
		Icon: nodes.NodeIcon{
			Brand: "telegram",
			Glyph: "image",
		},
		Inputs: []nodes.NodeHandle{
			{Key: "in"},
		},
		Outputs: []nodes.NodeHandle{
			{Key: "out"},
		},
		Categories:    []string{"Publishing"},
		SubCategories: []string{"Telegram"},
		Tags: []string{
			"chat",
			"messaging",
			"send",
			"media",
			"notification",
			"bot",
		},
		SupportedCredentials: []string{
			"telegram_api",
		},
		InputSchema: &gjs.Schema{
			Type: "object",
			Properties: map[string]*gjs.Schema{
				"chatId": {
					Type:        "string",
					Title:       "Chat ID",
					Description: "Telegram chat identifier where the media will be sent.",
				},
				"mediaUrls": {
					Type:        "array",
					Title:       "Media URLs",
					Description: "List of media URLs to send to the Telegram chat.",
					Items:       &gjs.Schema{Type: "string"},
				},
			},
			Required: []string{
				"chatId",
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
						Type:        "integer",
						Description: "Identifier of the sent Telegram message.",
					},
					"error": {
						Type:        "string",
						Description: "Error message when the media send failed.",
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
