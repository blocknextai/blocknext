package getplatformcredentialbyid

import (
	nodeEngineDomainCredentials "github.com/blocknextai/platform-api/internal/nodeengine/domain/credentials"
	platformDomainPlatformcredentials "github.com/blocknextai/platform-api/internal/platform/domain/platformcredentials"
)

func MapPlatformCredentialToResponse(
	platformCredential *platformDomainPlatformcredentials.PlatformCredential,
	credential nodeEngineDomainCredentials.CredentialManager,
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
