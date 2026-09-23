package oauth2

import (
	credentialOAuthApplicationOAuth2 "github.com/blocknextai/platform-api/internal/modules/credentialoauth/internal/application/oauth2"
	credentialsContract "github.com/blocknextai/platform-api/internal/modules/credentials/contract"
	nodeengineContract "github.com/blocknextai/platform-api/internal/modules/nodeengine/contract"
	platformContract "github.com/blocknextai/platform-api/internal/modules/platform/contract"
)

type Service struct {
	stateStore                  *credentialOAuthApplicationOAuth2.StateStore
	oauth2RedirectURL           string
	credentialService           credentialsContract.CredentialService
	nodeEngineCredentialService nodeengineContract.CredentialService
	platformCredentialService   platformContract.PlatformCredentialService
}

func NewService(
	stateStore *credentialOAuthApplicationOAuth2.StateStore,
	oauth2RedirectURL string,
	credentialService credentialsContract.CredentialService,
	nodeEngineCredentialService nodeengineContract.CredentialService,
	platformCredentialService platformContract.PlatformCredentialService,
) *Service {
	return &Service{
		stateStore:                  stateStore,
		oauth2RedirectURL:           oauth2RedirectURL,
		credentialService:           credentialService,
		nodeEngineCredentialService: nodeEngineCredentialService,
		platformCredentialService:   platformCredentialService,
	}
}
