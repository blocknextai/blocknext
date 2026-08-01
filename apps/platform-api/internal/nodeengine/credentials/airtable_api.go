package credentials

import (
	domain "github.com/blocknextai/platform-api/internal/nodeengine/domain/credentials"
	gjs "github.com/google/jsonschema-go/jsonschema"
)

func NewAirtableAPICredential() *domain.Credential {
	return &domain.Credential{
		ID:          "airtable_api",
		PlatformID:  "airtable_api",
		Name:        "Airtable",
		Description: "Airtable personal access token domain.",
		Icon: domain.CredentialIcon{
			Light: "airtable",
			Dark:  "airtable",
		},
		Schema: &gjs.Schema{
			Type: "object",
			Properties: map[string]*gjs.Schema{
				"accessToken": {
					Type:        "string",
					Title:       "Access Token",
					Description: "Airtable personal access token used to authorize API requests.",
					WriteOnly:   true,
				},
			},
			Required: []string{
				"accessToken",
			},
		},
		SupportedNodes: &[]string{
			"airtable_create_record",
			"airtable_list_records",
		},
	}
}
