package credentials

import (
	"github.com/blocknextai/go-packages/json"
	domain "github.com/blocknextai/platform-api/internal/nodeengine/domain/credentials"
	gjs "github.com/google/jsonschema-go/jsonschema"
)

const (
	facebookAuthorizationURL = "https://www.facebook.com/v23.0/dialog/oauth"
	facebookAccessTokenURL   = "https://graph.facebook.com/v23.0/oauth/access_token"
	facebookScope            = "pages_manage_posts pages_read_engagement business_management pages_show_list public_profile"
)

func NewFacebookOAuth2Credential(redirectURL string) *domain.Credential {
	return &domain.Credential{
		ID:          "facebook_oauth2",
		PlatformID:  "facebook_oauth2",
		Name:        "Facebook",
		Description: "Facebook OAuth2 credentials for managing pages and posts.",
		Icon: domain.CredentialIcon{
			Brand: "facebook",
		},
		Schema: &gjs.Schema{
			Type: "object",
			Properties: map[string]*gjs.Schema{
				"authUrl": {
					Type:    "string",
					Default: json.RawMessage(`"` + facebookAuthorizationURL + `"`),
					Extra: map[string]any{
						"hidden": true,
					},
				},
				"tokenUrl": {
					Type:    "string",
					Default: json.RawMessage(`"` + facebookAccessTokenURL + `"`),
					Extra: map[string]any{
						"hidden": true,
					},
				},
				"redirectUrl": {
					Type:        "string",
					Title:       "Redirect URL",
					Description: "OAuth2 redirect URL configured in your Facebook app.",
					Default:     json.RawMessage(`"` + redirectURL + `"`),
					ReadOnly:    true,
					Extra: map[string]any{
						"copyable": true,
					},
				},
				"scope": {
					Type:    "string",
					Default: json.RawMessage(`"` + facebookScope + `"`),
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
					Description: "Facebook app client ID.",
				},
				"clientSecret": {
					Type:        "string",
					Title:       "Client Secret",
					Description: "Facebook app client secret.",
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
			"facebook_publish_post",
			"facebook_publish_media_post",
			"facebook_publish_story",
		},
	}
}
