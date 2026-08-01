package sendemail

import (
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
	gjs "github.com/google/jsonschema-go/jsonschema"
)

type GmailSendEmailNode struct {
	nodes.Node
}

func NewGmailSendEmailNode(nodeID string) *GmailSendEmailNode {
	return &GmailSendEmailNode{
		Node: nodes.Node{
			ID:          nodeID,
			Version:     "0.0.1",
			Name:        "Gmail Send Email",
			Description: "Send an email through Gmail.",
			Icon: nodes.NodeIcon{
				Light: nodeID,
				Dark:  nodeID,
			},
			Categories:    []string{"Mailing"},
			SubCategories: []string{"Gmail"},
			Tags: []string{
				"google",
				"email",
				"mail",
				"send",
				"mailing",
				"notification",
			},
			SupportedCredentials: []string{
				"gmail_oauth2",
			},
			InputSchema: &gjs.Schema{
				Type: "object",
				Properties: map[string]*gjs.Schema{
					"to": {
						Type:        "string",
						Title:       "To",
						Description: "Recipient email address.",
						Format:      "email",
					},
					"subject": {
						Type:        "string",
						Title:       "Subject",
						Description: "Email subject line.",
					},
					"body": {
						Type:        "string",
						Title:       "Body",
						Description: "Plain text body of the email.",
					},
				},
				Required: []string{
					"to",
					"subject",
					"body",
				},
			},
			OutputSchema: &gjs.Schema{
				Type: "array",
				Items: &gjs.Schema{
					Type: "object",
					Properties: map[string]*gjs.Schema{
						"status": {
							Type:        "boolean",
							Description: "Whether the email was sent successfully.",
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
