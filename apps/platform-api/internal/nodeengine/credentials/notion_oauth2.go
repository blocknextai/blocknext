package credentials

import (
	"github.com/blocknextai/go-packages/json"
	domain "github.com/blocknextai/platform-api/internal/nodeengine/domain/credentials"
	gjs "github.com/google/jsonschema-go/jsonschema"
)

const (
	notionAuthorizationURL = "https://api.notion.com/v1/oauth/authorize"
	notionAccessTokenURL   = "https://api.notion.com/v1/oauth/token"
)

func NewNotionOAuth2Credential(redirectURL string) *domain.Credential {
	return &domain.Credential{
		ID:          "notion_oauth2",
		PlatformID:  "notion_oauth2",
		Name:        "Notion",
		Description: "Notion OAuth2 credentials for managing pages and databases.",
		Icon: domain.CredentialIcon{
			Light: "notion",
			Dark:  "notion",
		},
		Schema: &gjs.Schema{
			Type: "object",
			Properties: map[string]*gjs.Schema{
				"authUrl": {
					Type:    "string",
					Default: json.RawMessage(`"` + notionAuthorizationURL + `"`),
					Extra: map[string]any{
						"hidden": true,
					},
				},
				"tokenUrl": {
					Type:    "string",
					Default: json.RawMessage(`"` + notionAccessTokenURL + `"`),
					Extra: map[string]any{
						"hidden": true,
					},
				},
				"redirectUrl": {
					Type:        "string",
					Title:       "Redirect URL",
					Description: "OAuth2 redirect URL configured in your Notion integration.",
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
					Default: json.RawMessage(`"header"`),
					Extra: map[string]any{
						"hidden": true,
					},
				},
				"clientId": {
					Type:        "string",
					Title:       "Client ID",
					Description: "Notion integration client ID.",
				},
				"clientSecret": {
					Type:        "string",
					Title:       "Client Secret",
					Description: "Notion integration client secret.",
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
			"notion_create_page",
			"notion_update_page",
			"notion_search_pages",
			"notion_get_page",
		},
	}
}
