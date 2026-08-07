package credentials

import (
	"github.com/blocknextai/go-packages/json"
	domain "github.com/blocknextai/platform-api/internal/nodeengine/domain/credentials"
	gjs "github.com/google/jsonschema-go/jsonschema"
)

const (
	slackAuthorizationURL = "https://slack.com/oauth/v2/authorize"
	slackAccessTokenURL   = "https://slack.com/api/oauth.v2.access"
	slackScope            = "chat:write files:write"
	slackAuthQueryParams  = "user_scope=channels:read channels:write chat:write files:read files:write groups:read im:read mpim:read reactions:read reactions:write stars:read stars:write usergroups:write usergroups:read users.profile:read users.profile:write users:read"
)

func NewSlackOAuth2Credential(redirectURL string) *domain.Credential {
	return &domain.Credential{
		ID:          "slack_oauth2",
		PlatformID:  "slack_oauth2",
		Name:        "Slack",
		Description: "Slack OAuth2 credentials for sending messages and media.",
		Icon: domain.CredentialIcon{
			Brand: "slack",
		},
		Schema: &gjs.Schema{
			Type: "object",
			Properties: map[string]*gjs.Schema{
				"authUrl": {
					Type:    "string",
					Default: json.RawMessage(`"` + slackAuthorizationURL + `"`),
					Extra: map[string]any{
						"hidden": true,
					},
				},
				"tokenUrl": {
					Type:    "string",
					Default: json.RawMessage(`"` + slackAccessTokenURL + `"`),
					Extra: map[string]any{
						"hidden": true,
					},
				},
				"redirectUrl": {
					Type:        "string",
					Title:       "Redirect URL",
					Description: "OAuth2 redirect URL configured in your Slack app.",
					Default:     json.RawMessage(`"` + redirectURL + `"`),
					ReadOnly:    true,
					Extra: map[string]any{
						"copyable": true,
					},
				},
				"scope": {
					Type:    "string",
					Default: json.RawMessage(`"` + slackScope + `"`),
					Extra: map[string]any{
						"hidden": true,
					},
				},
				"authentication": {
					Type:    "string",
					Default: json.RawMessage(`"header"`),
					Extra: map[string]any{
						"hidden": true,
					},
				},
				"authQueryParameters": {
					Type:    "string",
					Default: json.RawMessage(`"` + slackAuthQueryParams + `"`),
					Extra: map[string]any{
						"hidden": true,
					},
				},
				"clientId": {
					Type:        "string",
					Title:       "Client ID",
					Description: "Slack app client ID.",
				},
				"clientSecret": {
					Type:        "string",
					Title:       "Client Secret",
					Description: "Slack app client secret.",
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
			"slack_send_message",
			"slack_send_media",
		},
	}
}
