package credentials

import (
	domain "github.com/blocknextai/platform-api/internal/nodeengine/domain/credentials"
	gjs "github.com/google/jsonschema-go/jsonschema"
)

func NewGeminiAPICredential() *domain.Credential {
	return &domain.Credential{
		ID:          "gemini_api",
		PlatformID:  "gemini_api",
		Name:        "Gemini",
		Description: "Google Gemini API domain.",
		Icon: domain.CredentialIcon{
			Light: "gemini",
			Dark:  "gemini",
		},
		Schema: &gjs.Schema{
			Type: "object",
			Properties: map[string]*gjs.Schema{
				"apiKey": {
					Type:        "string",
					Title:       "API Key",
					Description: "Google AI Studio API key for Gemini.",
					WriteOnly:   true,
				},
			},
			Required: []string{
				"apiKey",
			},
		},
		SupportedNodes: &[]string{
			"gemini_nano_banana",
			"gemini_chat",
		},
	}
}
