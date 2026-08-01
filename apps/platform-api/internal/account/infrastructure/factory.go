package infrastructure

import (
	"github.com/blocknextai/platform-api/internal/account/infrastructure/authproviders"
	"github.com/blocknextai/platform-api/internal/config"

	accountDomain "github.com/blocknextai/platform-api/internal/account/domain"
	accountDomainUserNonces "github.com/blocknextai/platform-api/internal/account/domain/usernonces"
	web3DomainCrypto "github.com/blocknextai/platform-api/internal/web3/domain/crypto"
)

func NewAuthProviderRegistry(options config.AuthOptions, signatureVerifier web3DomainCrypto.SignatureVerifier, userNonceRepository accountDomainUserNonces.UserNonceRepository) (*authproviders.AuthProviderRegistry, error) {
	authProviderRegistry := authproviders.NewAuthProviderRegistry()

	if options.Metamask.Enabled {
		metamaskProvider := authproviders.NewMetamaskAuthProvider(
			signatureVerifier,
			userNonceRepository,
		)
		authProviderRegistry.Register(accountDomain.AuthProviderMetamask, metamaskProvider)
	}

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
