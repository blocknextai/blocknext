package oauth2

type GrantType string

const (
	AuthorizationCodeGrant GrantType = "authorization_code"
	RefreshTokenGrant      GrantType = "refresh_token"
	PKCEGrant              GrantType = "pkce"
)

var (
	AllGrantTypes = map[GrantType]struct{}{
		AuthorizationCodeGrant: {},
		RefreshTokenGrant:      {},
		PKCEGrant:              {},
	}
)

func (gt GrantType) String() string {
	return string(gt)
}

func (gt GrantType) IsValid() bool {
	_, ok := AllGrantTypes[gt]
	return ok
}
