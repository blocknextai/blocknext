package sendmedia

import (
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
	gjs "github.com/google/jsonschema-go/jsonschema"
)

type WhatsAppSendMediaNode struct {
	nodes.Node
}

func NewWhatsAppSendMediaNode(nodeID string) *WhatsAppSendMediaNode {
	return &WhatsAppSendMediaNode{
		Node: nodes.Node{
			ID:          nodeID,
			Kind:        nodes.NodeKindAction,
			Version:     "0.0.1",
			Name:        "WhatsApp Send Media",
			Description: "Send a media message via WhatsApp.",
			Icon: nodes.NodeIcon{
				Brand: "whatsapp",
				Glyph: "image",
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
				"media",
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
					"mediaUrls": {
						Type:        "array",
						Title:       "Media URLs",
						Description: "List of media URLs to send to the recipient.",
						Items:       &gjs.Schema{Type: "string"},
					},
					"caption": {
						Type:        "string",
						Title:       "Caption",
						Description: "Optional caption to send alongside the media.",
					},
				},
				Required: []string{
					"phoneNumber",
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
							Description: "Identifier of the sent WhatsApp message.",
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
		},
	}
}
