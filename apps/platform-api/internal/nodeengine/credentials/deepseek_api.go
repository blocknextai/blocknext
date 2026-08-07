package credentials

import (
	domain "github.com/blocknextai/platform-api/internal/nodeengine/domain/credentials"
	gjs "github.com/google/jsonschema-go/jsonschema"
)

func NewDeepseekAPICredential() *domain.Credential {
	return &domain.Credential{
		ID:          "deepseek_api",
		PlatformID:  "deepseek_api",
		Name:        "DeepSeek",
		Description: "DeepSeek API credentials for AI chat models.",
		Icon: domain.CredentialIcon{
			Brand: "deepseek",
		},
		Schema: &gjs.Schema{
			Type: "object",
			Properties: map[string]*gjs.Schema{
				"apiKey": {
					Type:        "string",
					Title:       "API Key",
					Description: "DeepSeek API key from your account dashboard.",
					WriteOnly:   true,
				},
			},
			Required: []string{
				"apiKey",
			},
		},
		SupportedNodes: &[]string{
			"deepseek_chat",
		},
	}
}
