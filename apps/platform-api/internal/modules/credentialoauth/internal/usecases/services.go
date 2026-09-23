package usecases

import (
	credentialOAuthApplicationOAuth2 "github.com/blocknextai/platform-api/internal/modules/credentialoauth/internal/application/oauth2"
	oauth2UseCases "github.com/blocknextai/platform-api/internal/modules/credentialoauth/internal/usecases/oauth2"
	credentialsContract "github.com/blocknextai/platform-api/internal/modules/credentials/contract"
	nodeengineContract "github.com/blocknextai/platform-api/internal/modules/nodeengine/contract"
	platformContract "github.com/blocknextai/platform-api/internal/modules/platform/contract"
)

type Services struct {
	OAuth2 *oauth2UseCases.Service
}

type ServiceDependencies struct {
	OAuth2RedirectURL string

	CredentialService           credentialsContract.CredentialService
	NodeEngineCredentialService nodeengineContract.CredentialService
	PlatformCredentialService   platformContract.PlatformCredentialService
	StateStore                  *credentialOAuthApplicationOAuth2.StateStore
}

func NewServices(deps ServiceDependencies) *Services {
	return &Services{
		OAuth2: oauth2UseCases.NewService(deps.StateStore, deps.OAuth2RedirectURL, deps.CredentialService, deps.NodeEngineCredentialService, deps.PlatformCredentialService),
	}
}
