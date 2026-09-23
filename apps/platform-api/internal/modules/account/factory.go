package account

import (
	"github.com/blocknextai/platform-api/internal/config"
	"github.com/blocknextai/platform-api/internal/modules/account/internal/authproviders"

	accountDomain "github.com/blocknextai/platform-api/internal/modules/account/internal/domain"
	accountDomainUserNonces "github.com/blocknextai/platform-api/internal/modules/account/internal/domain/usernonces"
)

func NewAuthProviderRegistry(options config.AuthOptions, userNonceRepository accountDomainUserNonces.UserNonceRepository) (*authproviders.AuthProviderRegistry, error) {
	authProviderRegistry := authproviders.NewAuthProviderRegistry()

	if options.Google.Enabled {
		googleProvider := authproviders.NewGoogleAuthProvider(
			options.Google.ClientID,
			options.Google.ClientSecret,
			options.Google.RedirectURI,
			userNonceRepository,
		)
		authProviderRegistry.Register(accountDomain.AuthProviderGoogle, googleProvider)
	}

	if options.Github.Enabled {
		githubProvider := authproviders.NewGithubAuthProvider(
			options.Github.ClientID,
			options.Github.ClientSecret,
			options.Github.RedirectURI,
			userNonceRepository,
		)
		authProviderRegistry.Register(accountDomain.AuthProviderGithub, githubProvider)
	}

	if options.X.Enabled {
		xProvider := authproviders.NewXAuthProvider(
			options.X.ClientID,
			options.X.ClientSecret,
			options.X.RedirectURI,
			userNonceRepository,
		)
		authProviderRegistry.Register(accountDomain.AuthProviderX, xProvider)
	}

	if options.Facebook.Enabled {
		facebookProvider := authproviders.NewFacebookAuthProvider(
			options.Facebook.ClientID,
			options.Facebook.ClientSecret,
			options.Facebook.RedirectURI,
			userNonceRepository,
		)
		authProviderRegistry.Register(accountDomain.AuthProviderFacebook, facebookProvider)
	}

	return authProviderRegistry, nil
}
