package credentials

import (
	domain "github.com/blocknextai/platform-api/internal/nodeengine/domain/credentials"
	gjs "github.com/google/jsonschema-go/jsonschema"
)

func NewDiscordAPICredential() *domain.Credential {
	return &domain.Credential{
		ID:          "discord_api",
		PlatformID:  "discord_api",
		Name:        "Discord",
		Description: "Discord bot credentials for messaging and media.",
		Icon: domain.CredentialIcon{
			Brand: "discord",
		},
		Schema: &gjs.Schema{
			Type: "object",
			Properties: map[string]*gjs.Schema{
				"botToken": {
					Type:        "string",
					Title:       "Bot Token",
					Description: "Discord bot token from the Developer Portal.",
					WriteOnly:   true,
				},
			},
			Required: []string{
				"botToken",
			},
		},
		SupportedNodes: &[]string{
			"discord_send_message",
			"discord_send_media",
		},
	}
}
