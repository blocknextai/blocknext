package platformcredentials

import (
	platformDomain "github.com/blocknextai/platform-api/internal/platform/domain/platformcredentials"
	platformInfrastructurePlatformConfig "github.com/blocknextai/platform-api/internal/platform/infrastructure/platformconfig"
)

type PlatformCredentialService interface {
	GetAllPlatformCredentials() []*platformDomain.PlatformCredential
	GetPlatformCredential(id string) (*platformDomain.PlatformCredential, bool)
}

type platformCredentialService struct {
	platformConfigLoader *platformInfrastructurePlatformConfig.PlatformConfigLoader
}

func NewPlatformCredentialService(
	platformConfigLoader *platformInfrastructurePlatformConfig.PlatformConfigLoader,
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
