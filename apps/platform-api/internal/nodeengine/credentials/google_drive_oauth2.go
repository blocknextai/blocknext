package credentials

import (
	"github.com/blocknextai/go-packages/json"
	domain "github.com/blocknextai/platform-api/internal/nodeengine/domain/credentials"
	gjs "github.com/google/jsonschema-go/jsonschema"
)

const (
	googleDriveAuthorizationURL = "https://accounts.google.com/o/oauth2/v2/auth"
	googleDriveAccessTokenURL   = "https://oauth2.googleapis.com/token"
	googleDriveScope            = "https://www.googleapis.com/auth/drive https://www.googleapis.com/auth/drive.appdata https://www.googleapis.com/auth/drive.photos.readonly"
)

func NewGoogleDriveOAuth2Credential(redirectURL string) *domain.Credential {
	return &domain.Credential{
		ID:          "google_drive_oauth2",
		PlatformID:  "google_oauth2",
		Name:        "Google Drive",
		Description: "Google Drive OAuth2 credentials for managing files and folders.",
		Icon: domain.CredentialIcon{
			Brand: "google_drive",
		},
		Schema: &gjs.Schema{
			Type: "object",
			Properties: map[string]*gjs.Schema{
				"authUrl": {
					Type:    "string",
					Default: json.RawMessage(`"` + googleDriveAuthorizationURL + `"`),
					Extra: map[string]any{
						"hidden": true,
					},
				},
				"tokenUrl": {
					Type:    "string",
					Default: json.RawMessage(`"` + googleDriveAccessTokenURL + `"`),
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
					Default: json.RawMessage(`"` + googleDriveScope + `"`),
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
			"google_drive_create_folder",
			"google_drive_create_text_file",
			"google_drive_get_file",
			"google_drive_upload_file",
		},
	}
}
