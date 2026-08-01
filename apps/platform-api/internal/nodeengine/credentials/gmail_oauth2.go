package credentials

import (
	"github.com/blocknextai/go-packages/json"
	domain "github.com/blocknextai/platform-api/internal/nodeengine/domain/credentials"
	gjs "github.com/google/jsonschema-go/jsonschema"
)

const (
	gmailAuthorizationURL = "https://accounts.google.com/o/oauth2/v2/auth"
	gmailAccessTokenURL   = "https://oauth2.googleapis.com/token"
	gmailScope            = "https://www.googleapis.com/auth/gmail.labels https://www.googleapis.com/auth/gmail.addons.current.action.compose https://www.googleapis.com/auth/gmail.addons.current.message.action https://mail.google.com/ https://www.googleapis.com/auth/gmail.modify https://www.googleapis.com/auth/gmail.compose"
)

func NewGmailOAuth2Credential(redirectURL string) *domain.Credential {
	return &domain.Credential{
		ID:          "gmail_oauth2",
		PlatformID:  "google_oauth2",
		Name:        "Gmail",
		Description: "Gmail OAuth2 credentials for sending and reading email.",
		Icon: domain.CredentialIcon{
			Light: "gmail",
			Dark:  "gmail",
		},
		Schema: &gjs.Schema{
			Type: "object",
			Properties: map[string]*gjs.Schema{
				"authUrl": {
					Type:    "string",
					Default: json.RawMessage(`"` + gmailAuthorizationURL + `"`),
					Extra: map[string]any{
						"hidden": true,
					},
				},
				"tokenUrl": {
					Type:    "string",
					Default: json.RawMessage(`"` + gmailAccessTokenURL + `"`),
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
					Default: json.RawMessage(`"` + gmailScope + `"`),
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
			"gmail_organize_emails",
			"gmail_search_emails",
			"gmail_send_email",
		},
	}
}
