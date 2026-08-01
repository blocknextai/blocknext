package getcredentialbyid

import (
	commonDomainCredential "github.com/blocknextai/platform-api/internal/common/domain/credential"
	credentialsDomainCredentials "github.com/blocknextai/platform-api/internal/credentials/domain/credentials"
)

func MapCredentialToResponse(credential *credentialsDomainCredentials.Credential, processedData map[string]any) *GetCredentialByIDResponse {
	return &GetCredentialByIDResponse{
		ID:         credential.ID,
		Title:      credential.Title,
		Key:        credential.Key,
		UIKey:      commonDomainCredential.BuildUIKey(credential.OwnerType.String(), credential.ID.String()),
		Data:       processedData,
		SourceType: credential.SourceType,
		CreatedAt:  credential.CreatedAt,
		UpdatedAt:  credential.UpdatedAt,
	}
}
