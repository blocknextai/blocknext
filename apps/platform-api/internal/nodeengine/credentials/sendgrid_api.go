package credentials

import (
	domain "github.com/blocknextai/platform-api/internal/nodeengine/domain/credentials"
	gjs "github.com/google/jsonschema-go/jsonschema"
)

func NewSendgridAPICredential() *domain.Credential {
	return &domain.Credential{
		ID:          "sendgrid_api",
		PlatformID:  "sendgrid_api",
		Name:        "SendGrid",
		Description: "SendGrid API credentials for sending email.",
		Icon: domain.CredentialIcon{
			Brand: "sendgrid",
		},
		Schema: &gjs.Schema{
			Type: "object",
			Properties: map[string]*gjs.Schema{
				"apiKey": {
					Type:        "string",
					Title:       "API Key",
					Description: "SendGrid API key from your account dashboard.",
					WriteOnly:   true,
				},
			},
			Required: []string{
				"apiKey",
			},
		},
		SupportedNodes: &[]string{
			"sendgrid_send_email",
		},
	}
}
