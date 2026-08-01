package updatecredential

import (
	"context"
	"maps"

	"github.com/blocknextai/go-packages/database"
	"github.com/blocknextai/go-packages/secretmanager"
	commonDomainCredential "github.com/blocknextai/platform-api/internal/common/domain/credential"
	credentialsApplicationCredentials "github.com/blocknextai/platform-api/internal/credentials/application/credentials"
	credentialsDomainCredentials "github.com/blocknextai/platform-api/internal/credentials/domain/credentials"
	nodeEngineApplicationCredentials "github.com/blocknextai/platform-api/internal/nodeengine/application/credentials"
)

type Handler struct {
	credentialRepository        credentialsDomainCredentials.CredentialRepository
	secretManager               secretmanager.SecretManager
	transactionManager          database.TransactionManager
	credentialProcessor         nodeEngineApplicationCredentials.CredentialProcessor
	nodeEngineCredentialService nodeEngineApplicationCredentials.CredentialService
}

func New(
	credentialRepository credentialsDomainCredentials.CredentialRepository,
	secretManager secretmanager.SecretManager,
	transactionManager database.TransactionManager,
	credentialProcessor nodeEngineApplicationCredentials.CredentialProcessor,
	nodeEngineCredentialService nodeEngineApplicationCredentials.CredentialService,
) *Handler {
	return &Handler{
		credentialRepository:        credentialRepository,
		secretManager:               secretManager,
		transactionManager:          transactionManager,
		credentialProcessor:         credentialProcessor,
		nodeEngineCredentialService: nodeEngineCredentialService,
	}
}

func (h *Handler) Handle(ctx context.Context, request *UpdateCredentialCommand) (*UpdateCredentialResponse, error) {
	var response *UpdateCredentialResponse

	err := h.transactionManager.ExecuteInTransaction(ctx, func(txCtx context.Context) error {
		credential, err := h.credentialRepository.GetByIDAndOwner(txCtx, request.ID, request.OwnerType, request.OwnerID)
		if err != nil {
			return credentialsDomainCredentials.ErrCredentialNotFound
		}

		var existingData map[string]any
		err = h.secretManager.Decrypt(credential.Data, &existingData)
		if err != nil {
			return credentialsApplicationCredentials.ErrFailedToDecryptData.WithCause(err)
		}

		processedData := request.Data
		cred, exists := h.nodeEngineCredentialService.GetCredentialByID(request.Key)
		if exists {
			mergedData := make(map[string]any)
			maps.Copy(mergedData, existingData)
			schema := cred.GetSchema()
			for key, value := range request.Data {
				if !h.credentialProcessor.ShouldSkipField(schema, key, value) {
					mergedData[key] = value
				}
			}
			processedData = mergedData
		}

		encryptedData, err := h.secretManager.Encrypt(processedData)
		if err != nil {
			return credentialsApplicationCredentials.ErrFailedToEncryptData.WithCause(err)
		}

		credential, err = credential.Update(request.Key, request.Title, encryptedData)
		if err != nil {
			return err
		}

		err = h.credentialRepository.Update(txCtx, credential)
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
