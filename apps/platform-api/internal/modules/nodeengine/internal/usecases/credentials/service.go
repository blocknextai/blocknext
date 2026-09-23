package credentials

import (
	nodeEngineApplicationCredentials "github.com/blocknextai/platform-api/internal/modules/nodeengine/internal/application/credentials"
)

type Service struct {
	credentialService nodeEngineApplicationCredentials.CredentialService
}

func NewService(
	credentialService nodeEngineApplicationCredentials.CredentialService,
) *Service {
	return &Service{
		credentialService: credentialService,
	}
}
