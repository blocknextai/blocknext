package oauth2

type ErrorCode string

const (
	InvalidRequestError          ErrorCode = "invalid_request"
	InvalidClientError           ErrorCode = "invalid_client"
	InvalidGrantError            ErrorCode = "invalid_grant"
	UnauthorizedClientError      ErrorCode = "unauthorized_client"
	UnsupportedGrantTypeError    ErrorCode = "unsupported_grant_type"
	UnsupportedResponseTypeError ErrorCode = "unsupported_response_type"
	InvalidScopeError            ErrorCode = "invalid_scope"
	InvalidTargetError           ErrorCode = "invalid_target"
	AccessDeniedError            ErrorCode = "access_denied"
	InvalidRedirectURIError      ErrorCode = "invalid_redirect_uri"
	InvalidClientMetadataError   ErrorCode = "invalid_client_metadata"
	InvalidTokenError            ErrorCode = "invalid_token"
	InsufficientScopeError       ErrorCode = "insufficient_scope"
)

func (ec ErrorCode) String() string {
	return string(ec)
}
