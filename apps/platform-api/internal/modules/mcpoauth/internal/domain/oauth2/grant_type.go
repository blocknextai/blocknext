package oauth2

import (
	"slices"
)

type GrantType string

const (
	AuthorizationCodeGrant GrantType = "authorization_code"
	RefreshTokenGrant      GrantType = "refresh_token"
)

var (
	AllGrantTypes = map[GrantType]struct{}{
		AuthorizationCodeGrant: {},
		RefreshTokenGrant:      {},
	}
)

func (gt GrantType) String() string {
	return string(gt)
}

func (gt GrantType) IsValid() bool {
	_, ok := AllGrantTypes[gt]
	return ok
}

func SupportedGrantTypes(values []string) []GrantType {
	grantTypes := make([]GrantType, 0, len(values))
	for _, value := range values {
		grantType := GrantType(value)
		if grantType.IsValid() && !slices.Contains(grantTypes, grantType) {
			grantTypes = append(grantTypes, grantType)
		}
	}
	return grantTypes
}
