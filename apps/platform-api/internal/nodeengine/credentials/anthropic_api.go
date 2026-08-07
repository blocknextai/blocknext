package credentials

import (
	domain "github.com/blocknextai/platform-api/internal/nodeengine/domain/credentials"
	gjs "github.com/google/jsonschema-go/jsonschema"
)

func NewAnthropicAPICredential() *domain.Credential {
	return &domain.Credential{
		ID:          "anthropic_api",
		PlatformID:  "anthropic_api",
		Name:        "Anthropic",
		Description: "Anthropic API credentials for accessing Claude models.",
		Icon: domain.CredentialIcon{
			Brand: "anthropic",
		},
		Schema: &gjs.Schema{
			Type: "object",
			Properties: map[string]*gjs.Schema{
				"apiKey": {
					Type:        "string",
					Title:       "API Key",
					Description: "Anthropic API key from your console dashboard.",
					WriteOnly:   true,
				},
			},
			Required: []string{
				"apiKey",
			},
		},
		SupportedNodes: &[]string{
			"anthropic_chat",
		},
	}
}
