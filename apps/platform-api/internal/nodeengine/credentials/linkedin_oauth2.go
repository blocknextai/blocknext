package credentials

import (
	"github.com/blocknextai/go-packages/json"
	domain "github.com/blocknextai/platform-api/internal/nodeengine/domain/credentials"
	gjs "github.com/google/jsonschema-go/jsonschema"
)

const (
	linkedinAuthorizationURL = "https://www.linkedin.com/oauth/v2/authorization"
	linkedinAccessTokenURL   = "https://www.linkedin.com/oauth/v2/accessToken"
	linkedinScope            = "r_liteprofile r_emailaddress w_member_social"
)

func NewLinkedinOAuth2Credential(redirectURL string) *domain.Credential {
	return &domain.Credential{
		ID:          "linkedin_oauth2",
		PlatformID:  "linkedin_oauth2",
		Name:        "LinkedIn",
		Description: "LinkedIn OAuth2 credentials for publishing posts.",
		Icon: domain.CredentialIcon{
			Light: "linkedin",
			Dark:  "linkedin",
		},
		Schema: &gjs.Schema{
			Type: "object",
			Properties: map[string]*gjs.Schema{
				"authUrl": {
					Type:    "string",
					Default: json.RawMessage(`"` + linkedinAuthorizationURL + `"`),
					Extra: map[string]any{
						"hidden": true,
					},
				},
				"tokenUrl": {
					Type:    "string",
					Default: json.RawMessage(`"` + linkedinAccessTokenURL + `"`),
					Extra: map[string]any{
						"hidden": true,
					},
				},
				"redirectUrl": {
					Type:        "string",
					Title:       "Redirect URL",
					Description: "OAuth2 redirect URL configured in your LinkedIn app.",
					Default:     json.RawMessage(`"` + redirectURL + `"`),
					ReadOnly:    true,
					Extra: map[string]any{
						"copyable": true,
					},
				},
				"scope": {
					Type:    "string",
					Default: json.RawMessage(`"` + linkedinScope + `"`),
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
					Description: "LinkedIn app client ID.",
				},
				"clientSecret": {
					Type:        "string",
					Title:       "Client Secret",
					Description: "LinkedIn app client secret.",
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
			"linkedin_publish_post",
			"linkedin_publish_media_post",
		},
	}
}
