package getallplatformcredentials

import (
	"context"

	nodeEngineApplicationCredentials "github.com/blocknextai/platform-api/internal/nodeengine/application/credentials"
	platformApplicationPlatformCredentials "github.com/blocknextai/platform-api/internal/platform/application/platformcredentials"
)

type Handler struct {
	platformCredentialService platformApplicationPlatformCredentials.PlatformCredentialService
	credentialService         nodeEngineApplicationCredentials.CredentialService
}

func New(
	platformCredentialService platformApplicationPlatformCredentials.PlatformCredentialService,
	credentialService nodeEngineApplicationCredentials.CredentialService,
) *Handler {
	return &Handler{
		platformCredentialService: platformCredentialService,
		credentialService:         credentialService,
	}
}

func (h *Handler) Handle(ctx context.Context, request *GetAllPlatformCredentialsQuery) (*GetAllPlatformCredentialsResponse, error) {
	platformCredentials := h.platformCredentialService.GetAllPlatformCredentials()

	return new(MapPlatformCredentialsToResponse(platformCredentials, h.credentialService.GetCredentialByID)), nil
}
