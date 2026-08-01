package credentials

import (
	"github.com/blocknextai/go-packages/json"
	domain "github.com/blocknextai/platform-api/internal/nodeengine/domain/credentials"
	gjs "github.com/google/jsonschema-go/jsonschema"
)

const (
	instagramAuthorizationURL = "https://www.instagram.com/oauth/authorize"
	instagramAccessTokenURL   = "https://api.instagram.com/oauth/access_token"
	instagramScope            = "instagram_business_basic instagram_business_manage_messages instagram_business_manage_comments instagram_business_content_publish instagram_business_manage_insights"
)

func NewInstagramOAuth2Credential(redirectURL string) *domain.Credential {
	return &domain.Credential{
		ID:          "instagram_oauth2",
		PlatformID:  "instagram_oauth2",
		Name:        "Instagram",
		Description: "Instagram OAuth2 credentials for publishing posts and stories.",
		Icon: domain.CredentialIcon{
			Light: "instagram",
			Dark:  "instagram",
		},
		Schema: &gjs.Schema{
			Type: "object",
			Properties: map[string]*gjs.Schema{
				"authUrl": {
					Type:    "string",
					Default: json.RawMessage(`"` + instagramAuthorizationURL + `"`),
					Extra: map[string]any{
						"hidden": true,
					},
				},
				"tokenUrl": {
					Type:    "string",
					Default: json.RawMessage(`"` + instagramAccessTokenURL + `"`),
					Extra: map[string]any{
						"hidden": true,
					},
				},
				"redirectUrl": {
					Type:        "string",
					Title:       "Redirect URL",
					Description: "OAuth2 redirect URL configured in your Instagram app.",
					Default:     json.RawMessage(`"` + redirectURL + `"`),
					ReadOnly:    true,
					Extra: map[string]any{
						"copyable": true,
					},
				},
				"scope": {
					Type:    "string",
					Default: json.RawMessage(`"` + instagramScope + `"`),
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
					Description: "Instagram app client ID.",
				},
				"clientSecret": {
					Type:        "string",
					Title:       "Client Secret",
					Description: "Instagram app client secret.",
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
			"instagram_publish_post",
			"instagram_publish_story",
			"instagram_publish_reels",
		},
	}
}
