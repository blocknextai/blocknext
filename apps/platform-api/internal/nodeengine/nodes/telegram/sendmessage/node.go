package sendmessage

import (
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
	gjs "github.com/google/jsonschema-go/jsonschema"
)

type TelegramSendMessageNode struct {
	nodes.Node
}

func NewTelegramSendMessageNode(nodeID string) *TelegramSendMessageNode {
	return &TelegramSendMessageNode{
		ID:          nodeID,
		Kind:        nodes.NodeKindAction,
		Version:     "0.0.1",
		Name:        "Telegram Send Message",
		Description: "Send a text message to a Telegram chat.",
		Icon: nodes.NodeIcon{
			Brand: "telegram",
			Glyph: "chat",
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
			"message",
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
					Description: "Telegram chat identifier where the message will be sent.",
				},
				"text": {
					Type:        "string",
					Title:       "Text",
					Description: "Text content of the Telegram message.",
				},
			},
			Required: []string{
				"chatId",
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
						Description: "Whether the Telegram message was sent successfully.",
					},
					"messageId": {
						Type:        "integer",
						Description: "Identifier of the sent Telegram message.",
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
