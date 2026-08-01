package credentials

import (
	"time"

	commonDomainCredential "github.com/blocknextai/platform-api/internal/common/domain/credential"
	credentialsDomainCredentials "github.com/blocknextai/platform-api/internal/credentials/domain/credentials"
	"github.com/google/uuid"
)

type CredentialResponse struct {
	ID         uuid.UUID                               `json:"id"`
	Title      string                                  `json:"title"`
	Key        string                                  `json:"key"`
	UIKey      string                                  `json:"uiKey"`
	SourceType credentialsDomainCredentials.SourceType `json:"sourceType"`
	CreatedAt  time.Time                               `json:"createdAt"`
	UpdatedAt  time.Time                               `json:"updatedAt"`
}

func mapCredentialToResponse(c *credentialsDomainCredentials.Credential) CredentialResponse {
	return CredentialResponse{
		ID:         c.ID,
		Title:      c.Title,
		Key:        c.Key,
		UIKey:      commonDomainCredential.BuildUIKey(c.OwnerType.String(), c.ID.String()),
		SourceType: c.SourceType,
		CreatedAt:  c.CreatedAt,
		UpdatedAt:  c.UpdatedAt,
	}
}

func MapCredentialsToResponse(credentials []*credentialsDomainCredentials.Credential) []CredentialResponse {
	result := make([]CredentialResponse, 0, len(credentials))
	for _, c := range credentials {
		result = append(result, mapCredentialToResponse(c))
	}
	return result
}
