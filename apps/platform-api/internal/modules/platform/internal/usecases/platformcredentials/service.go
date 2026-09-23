package platformcredentials

import (
	nodeengineContract "github.com/blocknextai/platform-api/internal/modules/nodeengine/contract"
	platformApplicationPlatformCredentials "github.com/blocknextai/platform-api/internal/modules/platform/internal/application/platformcredentials"
)

type Service struct {
	platformCredentialService platformApplicationPlatformCredentials.PlatformCredentialService
	credentialService         nodeengineContract.CredentialService
}

func NewService(
	platformCredentialService platformApplicationPlatformCredentials.PlatformCredentialService,
	credentialService nodeengineContract.CredentialService,
) *Service {
	return &Service{
		platformCredentialService: platformCredentialService,
		credentialService:         credentialService,
	}
}
