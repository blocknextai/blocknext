package credentials

import (
	gjs "github.com/google/jsonschema-go/jsonschema"

	domain "github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/credentials"
)

func NewVeoAPICredential() *domain.Credential {
	return &domain.Credential{
		ID:          "veo_api",
		PlatformID:  "veo_api",
		Name:        "Veo",
		Description: "Google Veo API credentials for video generation.",
		Icon: domain.CredentialIcon{
			Brand: "veo",
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
