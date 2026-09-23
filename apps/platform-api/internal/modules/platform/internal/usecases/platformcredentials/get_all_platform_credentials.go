package platformcredentials

import (
	"context"

	nodeengineContract "github.com/blocknextai/platform-api/internal/modules/nodeengine/contract"
	platformDomainPlatformcredentials "github.com/blocknextai/platform-api/internal/modules/platform/internal/domain/platformcredentials"
)

type GetAllPlatformCredentialsQuery struct{}

type PlatformCredential struct {
	ID             string                            `json:"id"`
	Name           string                            `json:"name"`
	Description    string                            `json:"description"`
	Icon           nodeengineContract.CredentialIcon `json:"icon"`
	IsOAuth2       bool                              `json:"isOAuth2"`
	SupportedNodes *[]string                         `json:"supportedNodes"`
}

type GetAllPlatformCredentialsResponse = []PlatformCredential

func (s *Service) GetAllPlatformCredentials(ctx context.Context, request *GetAllPlatformCredentialsQuery) (*GetAllPlatformCredentialsResponse, error) {
	platformCredentials := s.platformCredentialService.GetAllPlatformCredentials()

	return new(MapPlatformCredentialsToResponse(platformCredentials, s.credentialService.GetCredentialByID)), nil
}

func MapPlatformCredentialsToResponse(
	platformCredentials []*platformDomainPlatformcredentials.PlatformCredential,
	credentialGetter func(id string) (nodeengineContract.CredentialManager, bool),
) GetAllPlatformCredentialsResponse {
	var result GetAllPlatformCredentialsResponse
	for _, platformCredential := range platformCredentials {
		credential, exists := credentialGetter(platformCredential.ID)
		if !exists {
			continue
		}

		result = append(result, MapPlatformCredentialToResponse(platformCredential, credential))
	}
	return result
}

func MapPlatformCredentialToResponse(
	platformCredential *platformDomainPlatformcredentials.PlatformCredential,
	credential nodeengineContract.CredentialManager,
) PlatformCredential {
	return PlatformCredential{
		ID:             platformCredential.ID,
		Name:           credential.GetName(),
		Description:    credential.GetDescription(),
		Icon:           credential.GetIcon(),
		IsOAuth2:       credential.GetIsOAuth2(),
		SupportedNodes: credential.GetSupportedNodes(),
	}
}
