package infrastructure

import (
	"github.com/blocknextai/platform-api/internal/common/application/cqrs"
	credentialOAuthApplicationOAuth2 "github.com/blocknextai/platform-api/internal/credentialoauth/application/oauth2"
	"github.com/blocknextai/platform-api/internal/credentialoauth/application/oauth2/authurl"
	"github.com/blocknextai/platform-api/internal/credentialoauth/application/oauth2/exchangecode"
	credentialsApplicationCredentials "github.com/blocknextai/platform-api/internal/credentials/application/credentials"
	nodeEngineApplicationCredentials "github.com/blocknextai/platform-api/internal/nodeengine/application/credentials"
	platformApplicationPlatformCredentials "github.com/blocknextai/platform-api/internal/platform/application/platformcredentials"
)

type Handlers struct {
	AuthURL      cqrs.Handler[*authurl.AuthURLCommand, *authurl.AuthURLResponse]
	ExchangeCode cqrs.Handler[*exchangecode.ExchangeCodeQuery, *exchangecode.ExchangeCodeResponse]
}

type RegisterInfrastructureDeps struct {
	OAuth2RedirectURL string

	CredentialService           credentialsApplicationCredentials.CredentialService
	NodeEngineCredentialService nodeEngineApplicationCredentials.CredentialService
	PlatformCredentialService   platformApplicationPlatformCredentials.PlatformCredentialService
	StateStore                  *credentialOAuthApplicationOAuth2.StateStore
}

func RegisterInfrastructure(deps RegisterInfrastructureDeps) *Handlers {
	return &Handlers{
		AuthURL:      cqrs.ValidationBehavior(authurl.New(deps.StateStore, deps.OAuth2RedirectURL, deps.CredentialService, deps.NodeEngineCredentialService, deps.PlatformCredentialService)),
		ExchangeCode: cqrs.ValidationBehavior(exchangecode.New(deps.StateStore, deps.OAuth2RedirectURL, deps.CredentialService, deps.NodeEngineCredentialService, deps.PlatformCredentialService)),
	}
}
