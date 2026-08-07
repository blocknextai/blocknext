package credentials

import (
	domain "github.com/blocknextai/platform-api/internal/nodeengine/domain/credentials"
	gjs "github.com/google/jsonschema-go/jsonschema"
)

func NewSunomusicAPICredential() *domain.Credential {
	return &domain.Credential{
		ID:          "sunomusic_api",
		PlatformID:  "sunomusic_api",
		Name:        "Suno Music",
		Description: "Suno Music API credentials for music and lyrics generation.",
		Icon: domain.CredentialIcon{
			Brand: "sunomusic",
		},
		Schema: &gjs.Schema{
			Type: "object",
			Properties: map[string]*gjs.Schema{
				"apiKey": {
					Type:        "string",
					Title:       "API Key",
					Description: "Suno API key from your account dashboard.",
					WriteOnly:   true,
				},
			},
			Required: []string{
				"apiKey",
			},
		},
		SupportedNodes: &[]string{
			"sunomusic_create_music",
			"sunomusic_create_lyrics",
		},
	}
}
