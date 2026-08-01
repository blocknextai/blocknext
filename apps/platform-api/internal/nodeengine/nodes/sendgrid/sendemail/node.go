package sendemail

import (
	"github.com/blocknextai/platform-api/internal/nodeengine/domain/nodes"
	gjs "github.com/google/jsonschema-go/jsonschema"
)

type SendgridSendEmailNode struct {
	nodes.Node
}

func NewSendgridSendEmailNode(nodeID string) *SendgridSendEmailNode {
	return &SendgridSendEmailNode{
		Node: nodes.Node{
			ID:          nodeID,
			Version:     "0.0.1",
			Name:        "SendGrid Send Email",
			Description: "Send an email through SendGrid.",
			Icon: nodes.NodeIcon{
				Light: nodeID,
				Dark:  nodeID,
			},
			Categories:    []string{"Mailing"},
			SubCategories: []string{"SendGrid"},
			Tags: []string{
				"email",
				"send",
				"transactional",
				"mailing",
				"notification",
			},
			SupportedCredentials: []string{
				"sendgrid_api",
			},
			InputSchema: &gjs.Schema{
				Type: "object",
				Properties: map[string]*gjs.Schema{
					"from": {
						Type:        "string",
						Title:       "From",
						Description: "Sender email address.",
						Format:      "email",
					},
					"to": {
						Type:        "array",
						Title:       "To",
						Description: "List of recipient email addresses.",
						Items:       &gjs.Schema{Type: "string", Format: "email"},
					},
					"subject": {
						Type:        "string",
						Title:       "Subject",
						Description: "Email subject line.",
					},
					"content": {
						Type:        "string",
						Title:       "Content",
						Description: "HTML body content of the email.",
					},
				},
				Required: []string{
					"from",
					"to",
					"subject",
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
