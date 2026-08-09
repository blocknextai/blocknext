package sendtextmessage

import (
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
	gjs "github.com/google/jsonschema-go/jsonschema"
)

type WhatsAppSendTextMessageNode struct {
	nodes.Node
}

func NewWhatsAppSendTextMessageNode(nodeID string) *WhatsAppSendTextMessageNode {
	return &WhatsAppSendTextMessageNode{
		Node: nodes.Node{
			ID:          nodeID,
			Kind:        nodes.NodeKindAction,
			Version:     "0.0.1",
			Name:        "WhatsApp Send Text Message",
			Description: "Send a text message via WhatsApp.",
			Icon: nodes.NodeIcon{
				Brand: "whatsapp",
				Glyph: "chat",
			},
			Inputs: []nodes.NodeHandle{
				{Key: "in"},
			},
			Outputs: []nodes.NodeHandle{
				{Key: "out"},
			},
			Categories:    []string{"Publishing"},
			SubCategories: []string{"WhatsApp"},
			Tags: []string{
				"chat",
				"messaging",
				"send",
				"message",
				"notification",
				"sms",
			},
			SupportedCredentials: []string{
				"whatsapp_api",
			},
			InputSchema: &gjs.Schema{
				Type: "object",
				Properties: map[string]*gjs.Schema{
					"phoneNumber": {
						Type:        "string",
						Title:       "Phone Number",
						Description: "Recipient phone number in international format.",
					},
					"message": {
						Type:        "string",
						Title:       "Message",
						Description: "Text message body to send.",
					},
				},
				Required: []string{
					"phoneNumber",
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
							Description: "Whether the message was sent successfully.",
						},
						"messageId": {
							Type:        "string",
							Description: "Identifier of the sent WhatsApp message.",
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
