package credentials

import (
	gjs "github.com/google/jsonschema-go/jsonschema"

	domain "github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/credentials"
)

func NewChatgptAPICredential() *domain.Credential {
	return &domain.Credential{
		ID:          "chatgpt_api",
		PlatformID:  "chatgpt_api",
		Name:        "ChatGPT",
		Description: "OpenAI API credentials for ChatGPT models.",
		Icon: domain.CredentialIcon{
			Brand: "chatgpt",
		},
		Schema: &gjs.Schema{
			Type: "object",
			Properties: map[string]*gjs.Schema{
				"apiKey": {
					Type:        "string",
					Title:       "API Key",
					Description: "OpenAI API key from your account dashboard.",
					WriteOnly:   true,
				},
			},
			Required: []string{
				"apiKey",
			},
		},
		SupportedNodes: &[]string{
			"chatgpt_chat",
		},
	}
}
