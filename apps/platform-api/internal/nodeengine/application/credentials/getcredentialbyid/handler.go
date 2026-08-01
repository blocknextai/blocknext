package getcredentialbyid

import (
	"context"

	nodeEngineApplicationCredentials "github.com/blocknextai/platform-api/internal/nodeengine/application/credentials"
	nodeEngineDomainCredentials "github.com/blocknextai/platform-api/internal/nodeengine/domain/credentials"
)

type Handler struct {
	credentialService nodeEngineApplicationCredentials.CredentialService
}

func New(
	credentialService nodeEngineApplicationCredentials.CredentialService,
) *Handler {
	return &Handler{
		credentialService: credentialService,
	}
}

func (h *Handler) Handle(_ context.Context, request *GetCredentialByIDQuery) (*GetCredentialByIDResponse, error) {
	credential, exists := h.credentialService.GetCredentialByID(request.ID)
	if !exists {
		return nil, nodeEngineDomainCredentials.ErrCredentialNotFound
	}

	return &credential, nil
}
