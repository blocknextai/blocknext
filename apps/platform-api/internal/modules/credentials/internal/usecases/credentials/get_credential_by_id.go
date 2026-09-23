package credentials

import (
	"context"
	"time"

	"github.com/google/uuid"

	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
	commonDomainCredential "github.com/blocknextai/platform-api/internal/common/domain/credential"
	credentialsApplicationCredentials "github.com/blocknextai/platform-api/internal/modules/credentials/internal/application/credentials"
	credentialsDomainCredentials "github.com/blocknextai/platform-api/internal/modules/credentials/internal/domain/credentials"
)

type GetCredentialByIDQuery struct {
	OwnerType    commonDomain.OwnerType
	OwnerID      uuid.UUID
	CredentialID uuid.UUID
}

type GetCredentialByIDResponse struct {
	ID         uuid.UUID                               `json:"id"`
	Title      string                                  `json:"title"`
	Key        string                                  `json:"key"`
	UIKey      string                                  `json:"uiKey"`
	Data       map[string]any                          `json:"data,omitempty"`
	SourceType credentialsDomainCredentials.SourceType `json:"sourceType"`
	CreatedAt  time.Time                               `json:"createdAt"`
	UpdatedAt  time.Time                               `json:"updatedAt"`
}

func (s *Service) GetCredentialByID(ctx context.Context, request *GetCredentialByIDQuery) (*GetCredentialByIDResponse, error) {
	credential, err := s.credentialRepository.GetByIDAndOwner(ctx, request.CredentialID, request.OwnerType, request.OwnerID)
	if err != nil {
		return nil, credentialsApplicationCredentials.ErrFailedToGetCredentialByID.WithCause(err)
	}

	var decryptedData map[string]any
	err = s.secretManager.Decrypt(credential.Data, &decryptedData)
	if err != nil {
		return nil, credentialsApplicationCredentials.ErrFailedToDecryptData.WithCause(err)
	}

	processedData := decryptedData
	cred, exists := s.nodeEngineCredentialService.GetCredentialByID(credential.Key)
	if exists {
		processedData = s.credentialProcessor.ProcessCredentialData(cred.GetSchema(), credential.Key, decryptedData)
	}

	return MapCredentialToResponse(credential, processedData), nil
}

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
