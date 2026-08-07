package credentials

import (
	"github.com/blocknextai/go-packages/json"
	domain "github.com/blocknextai/platform-api/internal/nodeengine/domain/credentials"
	gjs "github.com/google/jsonschema-go/jsonschema"
)

const (
	youtubeAuthorizationURL = "https://accounts.google.com/o/oauth2/v2/auth"
	youtubeAccessTokenURL   = "https://oauth2.googleapis.com/token"
	youtubeScope            = "https://www.googleapis.com/auth/youtube https://www.googleapis.com/auth/youtubepartner https://www.googleapis.com/auth/youtube.force-ssl https://www.googleapis.com/auth/youtube.upload https://www.googleapis.com/auth/youtubepartner-channel-audit"
)

func NewYoutubeOAuth2Credential(redirectURL string) *domain.Credential {
	return &domain.Credential{
		ID:          "youtube_oauth2",
		PlatformID:  "google_oauth2",
		Name:        "YouTube",
		Description: "YouTube OAuth2 credentials for uploading and managing videos.",
		Icon: domain.CredentialIcon{
			Brand: "youtube",
		},
		Schema: &gjs.Schema{
			Type: "object",
			Properties: map[string]*gjs.Schema{
				"authUrl": {
					Type:    "string",
					Default: json.RawMessage(`"` + youtubeAuthorizationURL + `"`),
					Extra: map[string]any{
						"hidden": true,
					},
				},
				"tokenUrl": {
					Type:    "string",
					Default: json.RawMessage(`"` + youtubeAccessTokenURL + `"`),
					Extra: map[string]any{
						"hidden": true,
					},
				},
				"redirectUrl": {
					Type:        "string",
					Title:       "Redirect URL",
					Description: "OAuth2 redirect URL configured in your Google Cloud project.",
					Default:     json.RawMessage(`"` + redirectURL + `"`),
					ReadOnly:    true,
					Extra: map[string]any{
						"copyable": true,
					},
				},
				"scope": {
					Type:    "string",
					Default: json.RawMessage(`"` + youtubeScope + `"`),
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
				"authQueryParameters": {
					Type:    "string",
					Default: json.RawMessage(`"access_type=offline&prompt=consent"`),
					Extra: map[string]any{
						"hidden": true,
					},
				},
				"clientId": {
					Type:        "string",
					Title:       "Client ID",
					Description: "Google OAuth2 client ID.",
				},
				"clientSecret": {
					Type:        "string",
					Title:       "Client Secret",
					Description: "Google OAuth2 client secret.",
					WriteOnly:   true,
				},
			},
			Required: []string{
				"authUrl",
				"tokenUrl",
				"redirectUrl",
				"scope",
				"authentication",
				"authQueryParameters",
				"clientId",
				"clientSecret",
			},
		},
		SupportedNodes: &[]string{
			"youtube_upload_video",
		},
	}
}
