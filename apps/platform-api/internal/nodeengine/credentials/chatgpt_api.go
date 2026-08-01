package credentials

import (
	domain "github.com/blocknextai/platform-api/internal/nodeengine/domain/credentials"
	gjs "github.com/google/jsonschema-go/jsonschema"
)

func NewChatgptAPICredential() *domain.Credential {
	return &domain.Credential{
		ID:          "chatgpt_api",
		PlatformID:  "chatgpt_api",
		Name:        "ChatGPT",
		Description: "OpenAI API credentials for ChatGPT models.",
		Icon: domain.CredentialIcon{
			Light: "chatgpt",
			Dark:  "chatgpt",
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
