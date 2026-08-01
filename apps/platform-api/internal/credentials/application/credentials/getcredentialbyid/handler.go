package getcredentialbyid

import (
	"context"

	"github.com/blocknextai/go-packages/secretmanager"
	credentialsApplicationCredentials "github.com/blocknextai/platform-api/internal/credentials/application/credentials"
	credentialsDomainCredentials "github.com/blocknextai/platform-api/internal/credentials/domain/credentials"
	nodeEngineApplicationCredentials "github.com/blocknextai/platform-api/internal/nodeengine/application/credentials"
)

type Handler struct {
	credentialRepository        credentialsDomainCredentials.CredentialRepository
	secretManager               secretmanager.SecretManager
	credentialProcessor         nodeEngineApplicationCredentials.CredentialProcessor
	nodeEngineCredentialService nodeEngineApplicationCredentials.CredentialService
}

func New(
	credentialRepository credentialsDomainCredentials.CredentialRepository,
	secretManager secretmanager.SecretManager,
	credentialProcessor nodeEngineApplicationCredentials.CredentialProcessor,
	nodeEngineCredentialService nodeEngineApplicationCredentials.CredentialService,
) *Handler {
	return &Handler{
		credentialRepository:        credentialRepository,
		secretManager:               secretManager,
		credentialProcessor:         credentialProcessor,
		nodeEngineCredentialService: nodeEngineCredentialService,
	}
}

func (h *Handler) Handle(ctx context.Context, request *GetCredentialByIDQuery) (*GetCredentialByIDResponse, error) {
	credential, err := h.credentialRepository.GetByIDAndOwner(ctx, request.CredentialID, request.OwnerType, request.OwnerID)
	if err != nil {
		return nil, credentialsApplicationCredentials.ErrFailedToGetCredentialByID.WithCause(err)
	}

	var decryptedData map[string]any
	err = h.secretManager.Decrypt(credential.Data, &decryptedData)
	if err != nil {
		return nil, credentialsApplicationCredentials.ErrFailedToDecryptData.WithCause(err)
	}

	processedData := decryptedData
	cred, exists := h.nodeEngineCredentialService.GetCredentialByID(credential.Key)
	if exists {
		processedData = h.credentialProcessor.ProcessCredentialData(cred.GetSchema(), credential.Key, decryptedData)
	}

	return MapCredentialToResponse(credential, processedData), nil
}
