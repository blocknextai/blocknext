package platformcredentials

import (
	"context"

	nodeengineContract "github.com/blocknextai/platform-api/internal/modules/nodeengine/contract"
	platformApplicationPlatformCredentials "github.com/blocknextai/platform-api/internal/modules/platform/internal/application/platformcredentials"
	platformDomainPlatformcredentials "github.com/blocknextai/platform-api/internal/modules/platform/internal/domain/platformcredentials"
)

type GetPlatformCredentialByIDQuery struct {
	ID string
}

type GetPlatformCredentialByIDResponse struct {
	ID             string                            `json:"id"`
	Name           string                            `json:"name"`
	Description    string                            `json:"description"`
	Icon           nodeengineContract.CredentialIcon `json:"icon"`
	IsOAuth2       bool                              `json:"isOAuth2"`
	SupportedNodes *[]string                         `json:"supportedNodes"`
}

func (s *Service) GetPlatformCredentialByID(ctx context.Context, request *GetPlatformCredentialByIDQuery) (*GetPlatformCredentialByIDResponse, error) {
	platformCredential, exists := s.platformCredentialService.GetPlatformCredential(request.ID)
	if !exists {
		return nil, platformApplicationPlatformCredentials.ErrPlatformCredentialNotFound
	}

	credential, exists := s.credentialService.GetCredentialByID(platformCredential.ID)
	if !exists {
		return nil, platformApplicationPlatformCredentials.ErrPlatformCredentialNotFound
	}

	return GetPlatformCredentialByIDMapPlatformCredentialToResponse(platformCredential, credential), nil
}

func GetPlatformCredentialByIDMapPlatformCredentialToResponse(
	platformCredential *platformDomainPlatformcredentials.PlatformCredential,
	credential nodeengineContract.CredentialManager,
) *GetPlatformCredentialByIDResponse {
	return &GetPlatformCredentialByIDResponse{
		ID:             platformCredential.ID,
		Name:           credential.GetName(),
		Description:    credential.GetDescription(),
		Icon:           credential.GetIcon(),
		IsOAuth2:       credential.GetIsOAuth2(),
		SupportedNodes: credential.GetSupportedNodes(),
	}
}
