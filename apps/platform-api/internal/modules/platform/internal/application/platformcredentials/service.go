package platformcredentials

import (
	platformDomain "github.com/blocknextai/platform-api/internal/modules/platform/internal/domain/platformcredentials"
	platformPlatformConfig "github.com/blocknextai/platform-api/internal/modules/platform/internal/platformconfig"
)

type PlatformCredentialService interface {
	GetAllPlatformCredentials() []*platformDomain.PlatformCredential
	GetPlatformCredential(id string) (*platformDomain.PlatformCredential, bool)
}

type platformCredentialService struct {
	platformConfigLoader *platformPlatformConfig.PlatformConfigLoader
}

func NewPlatformCredentialService(
	platformConfigLoader *platformPlatformConfig.PlatformConfigLoader,
) PlatformCredentialService {
	return &platformCredentialService{
		platformConfigLoader: platformConfigLoader,
	}
}

func (s *platformCredentialService) GetAllPlatformCredentials() []*platformDomain.PlatformCredential {
	return s.platformConfigLoader.GetAllPlatformCredentials()
}

func (s *platformCredentialService) GetPlatformCredential(id string) (*platformDomain.PlatformCredential, bool) {
	return s.platformConfigLoader.GetPlatformCredential(id)
}
