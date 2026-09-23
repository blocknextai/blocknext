package credentials

import (
	gjs "github.com/google/jsonschema-go/jsonschema"

	domain "github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/domain/credentials"
)

func NewPiapiAPICredential() *domain.Credential {
	return &domain.Credential{
		ID:          "piapi_api",
		PlatformID:  "piapi_api",
		Name:        "PiApi",
		Description: "PiAPI credentials for image, video, and audio generation.",
		Icon: domain.CredentialIcon{
			Brand: "piapi",
		},
		Schema: &gjs.Schema{
			Type: "object",
			Properties: map[string]*gjs.Schema{
				"apiKey": {
					Type:        "string",
					Title:       "API Key",
					Description: "PiAPI key from your account workspace.",
					WriteOnly:   true,
				},
			},
			Required: []string{
				"apiKey",
			},
		},
		SupportedNodes: &[]string{
			"piapi_image_gen",
			"piapi_video_gen",
			"piapi_audio_gen",
		},
	}
}
