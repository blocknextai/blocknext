package credentials

import (
	domain "github.com/blocknextai/platform-api/internal/nodeengine/domain/credentials"
	gjs "github.com/google/jsonschema-go/jsonschema"
)

func NewDeeplAPICredential() *domain.Credential {
	return &domain.Credential{
		ID:          "deepl_api",
		PlatformID:  "deepl_api",
		Name:        "DeepL",
		Description: "DeepL API credentials for translation services.",
		Icon: domain.CredentialIcon{
			Light: "deepl",
			Dark:  "deepl",
		},
		Schema: &gjs.Schema{
			Type: "object",
			Properties: map[string]*gjs.Schema{
				"apiKey": {
					Type:        "string",
					Title:       "API Key",
					Description: "DeepL API key from your account settings.",
					WriteOnly:   true,
				},
			},
			Required: []string{
				"apiKey",
			},
		},
		SupportedNodes: &[]string{
			"deepl_translate",
		},
	}
}
