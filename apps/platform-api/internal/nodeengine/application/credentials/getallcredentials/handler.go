package getallcredentials

import (
	"context"

	nodeEngineApplicationCredentials "github.com/blocknextai/platform-api/internal/nodeengine/application/credentials"
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

func (h *Handler) Handle(_ context.Context, _ *GetAllCredentialsQuery) (*GetAllCredentialsResponse, error) {
	return new(h.credentialService.GetAllCredentials()), nil
}
