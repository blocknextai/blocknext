package getplatformcredentialbyid

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

func (h *Handler) Handle(ctx context.Context, request *GetPlatformCredentialByIDQuery) (*GetPlatformCredentialByIDResponse, error) {
	platformCredential, exists := h.platformCredentialService.GetPlatformCredential(request.ID)
	if !exists {
		return nil, platformApplicationPlatformCredentials.ErrPlatformCredentialNotFound
	}

	credential, exists := h.credentialService.GetCredentialByID(platformCredential.ID)
	if !exists {
		return nil, platformApplicationPlatformCredentials.ErrPlatformCredentialNotFound
	}

	return MapPlatformCredentialToResponse(platformCredential, credential), nil
}
