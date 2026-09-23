package oauth2

type TokenEndpointAuthMethod string

const (
	NoneAuthentication              TokenEndpointAuthMethod = "none"
	ClientSecretBasicAuthentication TokenEndpointAuthMethod = "client_secret_basic"
	ClientSecretPostAuthentication  TokenEndpointAuthMethod = "client_secret_post"
)

var (
	AllTokenEndpointAuthMethods = map[TokenEndpointAuthMethod]struct{}{
		NoneAuthentication:              {},
		ClientSecretBasicAuthentication: {},
		ClientSecretPostAuthentication:  {},
	}
)

func (m TokenEndpointAuthMethod) String() string {
	return string(m)
}

func (m TokenEndpointAuthMethod) IsValid() bool {
	_, ok := AllTokenEndpointAuthMethods[m]
	return ok
}

func (m TokenEndpointAuthMethod) IsConfidential() bool {
	return m != NoneAuthentication
}
