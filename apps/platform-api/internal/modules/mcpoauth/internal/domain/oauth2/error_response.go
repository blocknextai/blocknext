package oauth2

type ErrorResponse struct {
	Error            ErrorCode `json:"error"`
	ErrorDescription string    `json:"error_description,omitempty"`
}
