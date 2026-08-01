package domain

type AuthProvider string

const (
	AuthProviderMetamask AuthProvider = "metamask"
	AuthProviderGoogle   AuthProvider = "google"
	AuthProviderX        AuthProvider = "x"
	AuthProviderFacebook AuthProvider = "facebook"
	AuthProviderGithub   AuthProvider = "github"
	AuthProviderEmail    AuthProvider = "email"
	AuthProviderPassword AuthProvider = "password"
)

var (
	AuthProviders = map[AuthProvider]AuthProviderCategory{
		AuthProviderMetamask: AuthProviderCategoryCrypto,
		AuthProviderGoogle:   AuthProviderCategoryOAuth,
		AuthProviderX:        AuthProviderCategoryOAuth,
		AuthProviderFacebook: AuthProviderCategoryOAuth,
		AuthProviderGithub:   AuthProviderCategoryOAuth,
		AuthProviderEmail:    AuthProviderCategoryEmail,
		AuthProviderPassword: AuthProviderCategoryPassword,
	}
)

func (a AuthProvider) String() string {
	return string(a)
}

func (a AuthProvider) IsValid() bool {
	_, ok := AuthProviders[a]
	return ok
}
