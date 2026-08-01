package credentials

import (
	domain "github.com/blocknextai/platform-api/internal/nodeengine/domain/credentials"
	gjs "github.com/google/jsonschema-go/jsonschema"
)

func NewWhatsappAPICredential() *domain.Credential {
	return &domain.Credential{
		ID:          "whatsapp_api",
		PlatformID:  "whatsapp_api",
		Name:        "WhatsApp",
		Description: "WhatsApp Business API credentials for sending messages.",
		Icon: domain.CredentialIcon{
			Light: "whatsapp",
			Dark:  "whatsapp",
		},
		Schema: &gjs.Schema{
			Type: "object",
			Properties: map[string]*gjs.Schema{
				"accessToken": {
					Type:        "string",
					Title:       "Access Token",
					Description: "WhatsApp Business permanent access token.",
					WriteOnly:   true,
				},
				"phoneNumberId": {
					Type:        "string",
					Title:       "Phone Number ID",
					Description: "WhatsApp Business phone number identifier.",
				},
			},
			Required: []string{
				"accessToken",
				"phoneNumberId",
			},
		},
		SupportedNodes: &[]string{
			"whatsapp_send_text_message",
			"whatsapp_send_template",
		},
	}
}
