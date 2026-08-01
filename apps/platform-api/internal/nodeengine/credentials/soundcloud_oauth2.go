package credentials

import (
	"github.com/blocknextai/go-packages/json"
	domain "github.com/blocknextai/platform-api/internal/nodeengine/domain/credentials"
	gjs "github.com/google/jsonschema-go/jsonschema"
)

const (
	soundcloudAuthorizationURL = "https://secure.soundcloud.com/authorize"
	soundcloudAccessTokenURL   = "https://secure.soundcloud.com/oauth/token"
)

func NewSoundcloudOAuth2Credential(redirectURL string) *domain.Credential {
	return &domain.Credential{
		ID:          "soundcloud_oauth2",
		PlatformID:  "soundcloud_oauth2",
		Name:        "SoundCloud",
		Description: "SoundCloud OAuth2 credentials for managing tracks and playlists.",
		Icon: domain.CredentialIcon{
			Light: "soundcloud",
			Dark:  "soundcloud",
		},
		Schema: &gjs.Schema{
			Type: "object",
			Properties: map[string]*gjs.Schema{
				"authUrl": {
					Type:    "string",
					Default: json.RawMessage(`"` + soundcloudAuthorizationURL + `"`),
					Extra: map[string]any{
						"hidden": true,
					},
				},
				"tokenUrl": {
					Type:    "string",
					Default: json.RawMessage(`"` + soundcloudAccessTokenURL + `"`),
					Extra: map[string]any{
						"hidden": true,
					},
				},
				"redirectUrl": {
					Type:        "string",
					Title:       "Redirect URL",
					Description: "OAuth2 redirect URL configured in your SoundCloud app.",
					Default:     json.RawMessage(`"` + redirectURL + `"`),
					ReadOnly:    true,
					Extra: map[string]any{
						"copyable": true,
					},
				},
				"scope": {
					Type:    "string",
					Default: json.RawMessage(`""`),
					Extra: map[string]any{
						"hidden": true,
					},
				},
				"authentication": {
					Type:    "string",
					Default: json.RawMessage(`"body"`),
					Extra: map[string]any{
						"hidden": true,
					},
				},
				"clientId": {
					Type:        "string",
					Title:       "Client ID",
					Description: "SoundCloud app client ID.",
				},
				"clientSecret": {
					Type:        "string",
					Title:       "Client Secret",
					Description: "SoundCloud app client secret.",
					WriteOnly:   true,
				},
			},
			Required: []string{
				"authUrl",
				"tokenUrl",
				"redirectUrl",
				"scope",
				"authentication",
				"clientId",
				"clientSecret",
			},
		},
		SupportedNodes: &[]string{
			"soundcloud_create_track",
			"soundcloud_create_playlist",
		},
	}
}
