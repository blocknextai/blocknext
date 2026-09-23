package oauth2

const (
	AuthorizationServerMetadataPath = "/.well-known/oauth-authorization-server"
	ProtectedResourceMetadataPath   = "/.well-known/oauth-protected-resource"

	AuthorizationEndpointPath = "/mcp-oauth/oauth2/authorize"
	TokenEndpointPath         = "/mcp-oauth/oauth2/token"
	RegistrationEndpointPath  = "/mcp-oauth/oauth2/register"
	RevocationEndpointPath    = "/mcp-oauth/oauth2/revoke"
)
