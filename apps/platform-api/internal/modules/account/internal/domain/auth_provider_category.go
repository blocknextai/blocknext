package domain

type AuthProviderCategory string

const (
	AuthProviderCategoryCrypto   AuthProviderCategory = "crypto"
	AuthProviderCategoryOAuth    AuthProviderCategory = "oauth"
	AuthProviderCategoryEmail    AuthProviderCategory = "email"
	AuthProviderCategoryPassword AuthProviderCategory = "password"
)

var (
	AllAuthProviderCategories = map[AuthProviderCategory]struct{}{
		AuthProviderCategoryCrypto:   {},
		AuthProviderCategoryOAuth:    {},
		AuthProviderCategoryEmail:    {},
		AuthProviderCategoryPassword: {},
	}
)

func (a AuthProvider) Category() AuthProviderCategory {
	return AuthProviders[a]
}

func (c AuthProviderCategory) String() string {
	return string(c)
}

func (c AuthProviderCategory) IsValid() bool {
	_, ok := AllAuthProviderCategories[c]
	return ok
}
