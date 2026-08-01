package createcredential

import (
	"context"

	"github.com/blocknextai/go-packages/database"
	"github.com/blocknextai/go-packages/secretmanager"
	commonDomainCredential "github.com/blocknextai/platform-api/internal/common/domain/credential"
	credentialsApplicationCredentials "github.com/blocknextai/platform-api/internal/credentials/application/credentials"
	credentialsDomainCredentials "github.com/blocknextai/platform-api/internal/credentials/domain/credentials"
	platformApplicationPlatformCredentials "github.com/blocknextai/platform-api/internal/platform/application/platformcredentials"
)

type Handler struct {
	credentialRepository      credentialsDomainCredentials.CredentialRepository
	secretManager             secretmanager.SecretManager
	transactionManager        database.TransactionManager
	platformCredentialService platformApplicationPlatformCredentials.PlatformCredentialService
}

func New(
	credentialRepository credentialsDomainCredentials.CredentialRepository,
	secretManager secretmanager.SecretManager,
	transactionManager database.TransactionManager,
	platformCredentialService platformApplicationPlatformCredentials.PlatformCredentialService,
) *Handler {
	return &Handler{
		credentialRepository:      credentialRepository,
		secretManager:             secretManager,
		transactionManager:        transactionManager,
		platformCredentialService: platformCredentialService,
	}
}

func (h *Handler) Handle(ctx context.Context, request *CreateCredentialCommand) (*CreateCredentialResponse, error) {
	if request.SourceType == credentialsDomainCredentials.SourceTypePlatform {
		if _, ok := h.platformCredentialService.GetPlatformCredential(request.Key); !ok {
			return nil, credentialsApplicationCredentials.ErrPlatformCredentialNotSupported
		}
	}

	var response *CreateCredentialResponse

	err := h.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		dataToEncrypt := request.Data
		if request.SourceType == credentialsDomainCredentials.SourceTypePlatform && dataToEncrypt == nil {
			dataToEncrypt = map[string]any{}
		}

		encryptedData, err := h.secretManager.Encrypt(dataToEncrypt)
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

		err = h.credentialRepository.Create(txCtx, credential)
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
