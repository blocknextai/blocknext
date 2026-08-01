package getallplatformcredentials

import (
	nodeEngineDomainCredentials "github.com/blocknextai/platform-api/internal/nodeengine/domain/credentials"
	platformDomainPlatformcredentials "github.com/blocknextai/platform-api/internal/platform/domain/platformcredentials"
)

func MapPlatformCredentialsToResponse(
	platformCredentials []*platformDomainPlatformcredentials.PlatformCredential,
	credentialGetter func(id string) (nodeEngineDomainCredentials.CredentialManager, bool),
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
	credential nodeEngineDomainCredentials.CredentialManager,
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
