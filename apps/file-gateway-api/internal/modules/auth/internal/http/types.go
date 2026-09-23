package http

const (
	HeaderServiceKey    = "X-Service-Key"
	HeaderAuthorization = "Authorization"
	BearerPrefix        = "Bearer "
)

type AuthTokenResponse struct {
	Token string `json:"token"`
}
