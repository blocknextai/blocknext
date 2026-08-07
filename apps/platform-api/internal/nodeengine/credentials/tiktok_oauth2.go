package credentials

import (
	"github.com/blocknextai/go-packages/json"
	domain "github.com/blocknextai/platform-api/internal/nodeengine/domain/credentials"
	gjs "github.com/google/jsonschema-go/jsonschema"
)

const (
	tiktokAuthorizationURL = "https://www.tiktok.com/v2/auth/authorize/"
	tiktokAccessTokenURL   = "https://open.tiktokapis.com/v2/oauth/token/"
	tiktokScope            = "user.info.basic,video.upload,video.publish"
)

func NewTiktokOAuth2Credential(redirectURL string) *domain.Credential {
	return &domain.Credential{
		ID:          "tiktok_oauth2",
		PlatformID:  "tiktok_oauth2",
		Name:        "TikTok",
		Description: "TikTok OAuth2 credentials for publishing videos.",
		Icon: domain.CredentialIcon{
			Brand: "tiktok",
		},
		Schema: &gjs.Schema{
			Type: "object",
			Properties: map[string]*gjs.Schema{
				"authUrl": {
					Type:    "string",
					Default: json.RawMessage(`"` + tiktokAuthorizationURL + `"`),
					Extra: map[string]any{
						"hidden": true,
					},
				},
				"tokenUrl": {
					Type:    "string",
					Default: json.RawMessage(`"` + tiktokAccessTokenURL + `"`),
					Extra: map[string]any{
						"hidden": true,
					},
				},
				"redirectUrl": {
					Type:        "string",
					Title:       "Redirect URL",
					Description: "OAuth2 redirect URL configured in your TikTok app.",
					Default:     json.RawMessage(`"` + redirectURL + `"`),
					ReadOnly:    true,
					Extra: map[string]any{
						"copyable": true,
					},
				},
				"scope": {
					Type:    "string",
					Default: json.RawMessage(`"` + tiktokScope + `"`),
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
				"clientKey": {
					Type:        "string",
					Title:       "Client Key",
					Description: "TikTok app client key.",
				},
				"clientSecret": {
					Type:        "string",
					Title:       "Client Secret",
					Description: "TikTok app client secret.",
					WriteOnly:   true,
				},
			},
			Required: []string{
				"authUrl",
				"tokenUrl",
				"redirectUrl",
				"scope",
				"authentication",
				"clientKey",
				"clientSecret",
			},
		},
		SupportedNodes: &[]string{
			"tiktok_publish_post",
		},
	}
}
