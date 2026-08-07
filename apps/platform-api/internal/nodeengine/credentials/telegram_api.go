package credentials

import (
	domain "github.com/blocknextai/platform-api/internal/nodeengine/domain/credentials"
	gjs "github.com/google/jsonschema-go/jsonschema"
)

func NewTelegramAPICredential() *domain.Credential {
	return &domain.Credential{
		ID:          "telegram_api",
		PlatformID:  "telegram_api",
		Name:        "Telegram",
		Description: "Telegram bot credentials for messaging.",
		Icon: domain.CredentialIcon{
			Brand: "telegram",
		},
		Schema: &gjs.Schema{
			Type: "object",
			Properties: map[string]*gjs.Schema{
				"botToken": {
					Type:        "string",
					Title:       "Bot Token",
					Description: "Telegram bot token issued by BotFather.",
					WriteOnly:   true,
				},
			},
			Required: []string{
				"botToken",
			},
		},
		SupportedNodes: &[]string{
			"telegram_send_message",
			"telegram_send_media",
		},
	}
}
