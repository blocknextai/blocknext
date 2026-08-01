package credentials

import (
	domain "github.com/blocknextai/platform-api/internal/nodeengine/domain/credentials"
	gjs "github.com/google/jsonschema-go/jsonschema"
)

func NewElevenlabsAPICredential() *domain.Credential {
	return &domain.Credential{
		ID:          "elevenlabs_api",
		PlatformID:  "elevenlabs_api",
		Name:        "ElevenLabs",
		Description: "ElevenLabs API credentials for text-to-speech generation.",
		Icon: domain.CredentialIcon{
			Light: "elevenlabs",
			Dark:  "elevenlabs",
		},
		Schema: &gjs.Schema{
			Type: "object",
			Properties: map[string]*gjs.Schema{
				"apiKey": {
					Type:        "string",
					Title:       "API Key",
					Description: "ElevenLabs API key from your account profile.",
					WriteOnly:   true,
				},
			},
			Required: []string{
				"apiKey",
			},
		},
		SupportedNodes: &[]string{
			"elevenlabs_text_to_speech",
		},
	}
}
