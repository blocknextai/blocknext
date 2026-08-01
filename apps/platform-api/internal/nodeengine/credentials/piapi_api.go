package credentials

import (
	domain "github.com/blocknextai/platform-api/internal/nodeengine/domain/credentials"
	gjs "github.com/google/jsonschema-go/jsonschema"
)

func NewPiapiAPICredential() *domain.Credential {
	return &domain.Credential{
		ID:          "piapi_api",
		PlatformID:  "piapi_api",
		Name:        "PiApi",
		Description: "PiAPI credentials for image, video, and audio generation.",
		Icon: domain.CredentialIcon{
			Light: "piapi",
			Dark:  "piapi",
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
