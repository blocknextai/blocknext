package oauth2

type AuthenticationMethod string

const (
	HeaderAuthentication AuthenticationMethod = "header"
	BodyAuthentication   AuthenticationMethod = "body"
)

var (
	AllAuthenticationMethods = map[AuthenticationMethod]struct{}{
		HeaderAuthentication: {},
		BodyAuthentication:   {},
	}
)

func (am AuthenticationMethod) String() string {
	return string(am)
}

func (am AuthenticationMethod) IsValid() bool {
	_, ok := AllAuthenticationMethods[am]
	return ok
}
