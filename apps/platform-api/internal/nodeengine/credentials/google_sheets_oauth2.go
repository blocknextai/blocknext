package credentials

import (
	"github.com/blocknextai/go-packages/json"
	domain "github.com/blocknextai/platform-api/internal/nodeengine/domain/credentials"
	gjs "github.com/google/jsonschema-go/jsonschema"
)

const (
	googleSheetsAuthorizationURL = "https://accounts.google.com/o/oauth2/v2/auth"
	googleSheetsAccessTokenURL   = "https://oauth2.googleapis.com/token"
	googleSheetsScope            = "https://www.googleapis.com/auth/spreadsheets https://www.googleapis.com/auth/drive https://www.googleapis.com/auth/drive.file https://www.googleapis.com/auth/drive.metadata"
)

func NewGoogleSheetsOAuth2Credential(redirectURL string) *domain.Credential {
	return &domain.Credential{
		ID:          "google_sheets_oauth2",
		PlatformID:  "google_oauth2",
		Name:        "Google Sheets",
		Description: "Google Sheets OAuth2 credentials for managing spreadsheets.",
		Icon: domain.CredentialIcon{
			Brand: "google_sheets",
		},
		Schema: &gjs.Schema{
			Type: "object",
			Properties: map[string]*gjs.Schema{
				"authUrl": {
					Type:    "string",
					Default: json.RawMessage(`"` + googleSheetsAuthorizationURL + `"`),
					Extra: map[string]any{
						"hidden": true,
					},
				},
				"tokenUrl": {
					Type:    "string",
					Default: json.RawMessage(`"` + googleSheetsAccessTokenURL + `"`),
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
					Default: json.RawMessage(`"` + googleSheetsScope + `"`),
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
			"sheets_create_spreadsheet",
			"sheets_read_spreadsheet",
			"sheets_delete_spreadsheet",
			"sheets_add_data_to_spreadsheet",
			"sheets_update",
		},
	}
}
