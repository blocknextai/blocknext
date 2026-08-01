package sendtemplate

import (
	"github.com/blocknextai/go-packages/json"
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
	gjs "github.com/google/jsonschema-go/jsonschema"
)

type WhatsAppSendTemplateNode struct {
	nodes.Node
}

func NewWhatsAppSendTemplateNode(nodeID string) *WhatsAppSendTemplateNode {
	return &WhatsAppSendTemplateNode{
		Node: nodes.Node{
			ID:          nodeID,
			Version:     "0.0.1",
			Name:        "WhatsApp Send Template",
			Description: "Send a template message via WhatsApp.",
			Icon: nodes.NodeIcon{
				Light: nodeID,
				Dark:  nodeID,
			},
			Categories:    []string{"Publishing"},
			SubCategories: []string{"WhatsApp"},
			Tags: []string{
				"chat",
				"messaging",
				"send",
				"template",
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
					"templateName": {
						Type:        "string",
						Title:       "Template Name",
						Description: "Name of the WhatsApp template to send.",
					},
					"languageCode": {
						Type:        "string",
						Title:       "Language Code",
						Description: "Locale code for the template (e.g. en_US).",
						Default:     json.RawMessage(`"en_US"`),
					},
				},
				Required: []string{
					"phoneNumber",
					"templateName",
					"languageCode",
				},
			},
			OutputSchema: &gjs.Schema{
				Type: "array",
				Items: &gjs.Schema{
					Type: "object",
					Properties: map[string]*gjs.Schema{
						"status": {
							Type:        "boolean",
							Description: "Whether the template message was sent successfully.",
						},
						"messageId": {
							Type:        "string",
							Description: "Identifier of the sent WhatsApp template message.",
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
