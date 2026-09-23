package credentials

import (
	"context"
	"maps"
	"strings"

	"github.com/google/uuid"

	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
	commonDomainCredential "github.com/blocknextai/platform-api/internal/common/domain/credential"
	credentialsApplicationCredentials "github.com/blocknextai/platform-api/internal/modules/credentials/internal/application/credentials"
	credentialsDomainCredentials "github.com/blocknextai/platform-api/internal/modules/credentials/internal/domain/credentials"
)

type UpdateCredentialCommand struct {
	ID        uuid.UUID
	OwnerType commonDomain.OwnerType
	OwnerID   uuid.UUID
	Key       string
	Title     string
	Data      map[string]any
}

const (
	UpdateCredentialMaxKeyLength   = 255
	UpdateCredentialMaxTitleLength = 255
)

func (c *UpdateCredentialCommand) Validate() error {
	if strings.TrimSpace(c.Key) == "" {
		return credentialsDomainCredentials.ErrInvalidKey
	}

	if len(strings.TrimSpace(c.Key)) > UpdateCredentialMaxKeyLength {
		return credentialsDomainCredentials.ErrKeyTooLong
	}

	if strings.TrimSpace(c.Title) == "" {
		return credentialsDomainCredentials.ErrInvalidTitle
	}

	if len(strings.TrimSpace(c.Title)) > UpdateCredentialMaxTitleLength {
		return credentialsDomainCredentials.ErrTitleTooLong
	}

	if len(c.Data) == 0 {
		return credentialsDomainCredentials.ErrInvalidData
	}

	return nil
}

type UpdateCredentialResponse struct {
	ID    uuid.UUID `json:"id"`
	Title string    `json:"title"`
	Key   string    `json:"key"`
	UIKey string    `json:"uiKey"`
}

func (s *Service) UpdateCredential(ctx context.Context, request *UpdateCredentialCommand) (*UpdateCredentialResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, err
	}

	var response *UpdateCredentialResponse

	err := s.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		credential, err := s.credentialRepository.GetByIDAndOwner(txCtx, request.ID, request.OwnerType, request.OwnerID)
		if err != nil {
			return credentialsDomainCredentials.ErrCredentialNotFound
		}

		var existingData map[string]any
		err = s.secretManager.Decrypt(credential.Data, &existingData)
		if err != nil {
			return credentialsApplicationCredentials.ErrFailedToDecryptData.WithCause(err)
		}

		processedData := request.Data
		cred, exists := s.nodeEngineCredentialService.GetCredentialByID(request.Key)
		if exists {
			mergedData := make(map[string]any)
			maps.Copy(mergedData, existingData)
			schema := cred.GetSchema()
			for key, value := range request.Data {
				if !s.credentialProcessor.ShouldSkipField(schema, key, value) {
					mergedData[key] = value
				}
			}
			processedData = mergedData
		}

		encryptedData, err := s.secretManager.Encrypt(processedData)
		if err != nil {
			return credentialsApplicationCredentials.ErrFailedToEncryptData.WithCause(err)
		}

		credential, err = credential.Update(request.Key, request.Title, encryptedData)
		if err != nil {
			return err
		}

		err = s.credentialRepository.Update(txCtx, credential)
		if err != nil {
			return credentialsApplicationCredentials.ErrFailedToUpdateCredential.WithCause(err)
		}

		uiKey := commonDomainCredential.BuildUIKey(credential.OwnerType.String(), credential.ID.String())

		response = &UpdateCredentialResponse{
			ID:    credential.ID,
			Title: credential.Title,
			Key:   credential.Key,
			UIKey: uiKey,
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return response, nil
}
