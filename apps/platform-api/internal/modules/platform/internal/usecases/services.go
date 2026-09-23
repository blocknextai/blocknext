package usecases

import (
	nodeengineContract "github.com/blocknextai/platform-api/internal/modules/nodeengine/contract"
	platformApplicationPlatformCredentials "github.com/blocknextai/platform-api/internal/modules/platform/internal/application/platformcredentials"
	featuresUseCases "github.com/blocknextai/platform-api/internal/modules/platform/internal/usecases/features"
	platformcredentialsUseCases "github.com/blocknextai/platform-api/internal/modules/platform/internal/usecases/platformcredentials"
)

type Services struct {
	Features            *featuresUseCases.Service
	PlatformCredentials *platformcredentialsUseCases.Service
}

type ServiceDependencies struct {
	FunctionCallingEnabled     bool
	WorkflowsGenerationEnabled bool

	PlatformCredentialService platformApplicationPlatformCredentials.PlatformCredentialService
	CredentialService         nodeengineContract.CredentialService
}

func NewServices(deps ServiceDependencies) *Services {
	return &Services{
		Features:            featuresUseCases.NewService(deps.FunctionCallingEnabled, deps.WorkflowsGenerationEnabled),
		PlatformCredentials: platformcredentialsUseCases.NewService(deps.PlatformCredentialService, deps.CredentialService),
	}
}
