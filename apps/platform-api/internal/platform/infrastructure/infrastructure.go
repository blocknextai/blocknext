package infrastructure

import (
	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
	nodeEngineApplicationCredentials "github.com/blocknextai/platform-api/internal/nodeengine/application/credentials"
	"github.com/blocknextai/platform-api/internal/platform/application/features/getfeatures"
	platformApplicationPlatformCredentials "github.com/blocknextai/platform-api/internal/platform/application/platformcredentials"
	"github.com/blocknextai/platform-api/internal/platform/application/platformcredentials/getallplatformcredentials"
	"github.com/blocknextai/platform-api/internal/platform/application/platformcredentials/getplatformcredentialbyid"
)

type Handlers struct {
	GetAllPlatformCredentials cqrs.Handler[*getallplatformcredentials.GetAllPlatformCredentialsQuery, *getallplatformcredentials.GetAllPlatformCredentialsResponse]
	GetPlatformCredentialByID cqrs.Handler[*getplatformcredentialbyid.GetPlatformCredentialByIDQuery, *getplatformcredentialbyid.GetPlatformCredentialByIDResponse]
	GetFeatures               cqrs.Handler[*getfeatures.GetFeaturesQuery, *getfeatures.GetFeaturesResponse]
}

type RegisterInfrastructureDeps struct {
	FunctionCallingEnabled     bool
	WorkflowsGenerationEnabled bool

	PlatformCredentialService platformApplicationPlatformCredentials.PlatformCredentialService
	CredentialService         nodeEngineApplicationCredentials.CredentialService
}

func RegisterInfrastructure(deps RegisterInfrastructureDeps) *Handlers {
	return &Handlers{
		GetAllPlatformCredentials: cqrs.ValidationBehavior(getallplatformcredentials.New(deps.PlatformCredentialService, deps.CredentialService)),
		GetPlatformCredentialByID: cqrs.ValidationBehavior(getplatformcredentialbyid.New(deps.PlatformCredentialService, deps.CredentialService)),
		GetFeatures:               cqrs.ValidationBehavior(getfeatures.New(deps.FunctionCallingEnabled, deps.WorkflowsGenerationEnabled)),
	}
}
