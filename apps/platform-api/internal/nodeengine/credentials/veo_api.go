package credentials

import (
	domain "github.com/blocknextai/platform-api/internal/nodeengine/domain/credentials"
	gjs "github.com/google/jsonschema-go/jsonschema"
)

func NewVeoAPICredential() *domain.Credential {
	return &domain.Credential{
		ID:          "veo_api",
		PlatformID:  "veo_api",
		Name:        "Veo",
		Description: "Google Veo API credentials for video generation.",
		Icon: domain.CredentialIcon{
			Light: "veo",
			Dark:  "veo",
		},
		Schema: &gjs.Schema{
			Type: "object",
			Properties: map[string]*gjs.Schema{
				"apiKey": {
					Type:        "string",
					Title:       "API Key",
					Description: "Google AI Studio API key for Veo.",
					WriteOnly:   true,
				},
			},
			Required: []string{
				"apiKey",
			},
		},
		SupportedNodes: &[]string{
			"veo",
		},
	}
}
