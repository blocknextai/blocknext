package oauth2

const (
	HeaderBearerMethod = "header"
)

type ProtectedResourceMetadata struct {
	Resource               string
	AuthorizationServers   []string
	ScopesSupported        []string
	BearerMethodsSupported []string
}
