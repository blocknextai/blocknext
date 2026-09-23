package credentials

import (
	"context"
	"strings"

	"github.com/google/uuid"

	commonDomain "github.com/blocknextai/platform-api/internal/common/domain"
	commonDomainCredential "github.com/blocknextai/platform-api/internal/common/domain/credential"
	credentialsApplicationCredentials "github.com/blocknextai/platform-api/internal/modules/credentials/internal/application/credentials"
	credentialsDomainCredentials "github.com/blocknextai/platform-api/internal/modules/credentials/internal/domain/credentials"
)

type CreateCredentialCommand struct {
	OwnerType  commonDomain.OwnerType
	OwnerID    uuid.UUID
	SourceType credentialsDomainCredentials.SourceType
	Key        string
	Title      string
	Data       map[string]any
}

const (
	MaxKeyLength   = 255
	MaxTitleLength = 255
)

func (c *CreateCredentialCommand) Validate() error {
	if !c.SourceType.IsValid() {
		return credentialsDomainCredentials.ErrInvalidSourceType
	}

	if strings.TrimSpace(c.Key) == "" {
		return credentialsDomainCredentials.ErrInvalidKey
	}

	if len(strings.TrimSpace(c.Key)) > MaxKeyLength {
		return credentialsDomainCredentials.ErrKeyTooLong
	}

	if strings.TrimSpace(c.Title) == "" {
		return credentialsDomainCredentials.ErrInvalidTitle
	}

	if len(strings.TrimSpace(c.Title)) > MaxTitleLength {
		return credentialsDomainCredentials.ErrTitleTooLong
	}

	if c.SourceType == credentialsDomainCredentials.SourceTypeOwner {
		if len(c.Data) == 0 {
			return credentialsDomainCredentials.ErrInvalidData
		}
	}

	return nil
}

type CreateCredentialResponse struct {
	ID    uuid.UUID `json:"id"`
	Title string    `json:"title"`
	Key   string    `json:"key"`
	UIKey string    `json:"uiKey"`
}

func (s *Service) CreateCredential(ctx context.Context, request *CreateCredentialCommand) (*CreateCredentialResponse, error) {
	if err := request.Validate(); err != nil {
		return nil, err
	}

	if request.SourceType == credentialsDomainCredentials.SourceTypePlatform {
		if _, ok := s.platformCredentialService.GetPlatformCredential(request.Key); !ok {
			return nil, credentialsApplicationCredentials.ErrPlatformCredentialNotSupported
		}
	}

	var response *CreateCredentialResponse

	err := s.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		dataToEncrypt := request.Data
		if request.SourceType == credentialsDomainCredentials.SourceTypePlatform && dataToEncrypt == nil {
			dataToEncrypt = map[string]any{}
		}

		encryptedData, err := s.secretManager.Encrypt(dataToEncrypt)
		if err != nil {
			return credentialsApplicationCredentials.ErrFailedToEncryptData.WithCause(err)
		}

		credential, err := credentialsDomainCredentials.New(
			request.OwnerType,
			request.OwnerID,
			request.SourceType,
			request.Key,
			request.Title,
			encryptedData,
		)
		if err != nil {
			return err
		}

		err = s.credentialRepository.Create(txCtx, credential)
		if err != nil {
			return credentialsApplicationCredentials.ErrFailedToCreateCredential.WithCause(err)
		}

		uiKey := commonDomainCredential.BuildUIKey(request.OwnerType.String(), credential.ID.String())

		response = &CreateCredentialResponse{
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
