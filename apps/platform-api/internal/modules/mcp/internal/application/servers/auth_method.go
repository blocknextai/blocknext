package servers

type AuthMethod string

const (
	AuthMethodAPIKey AuthMethod = "apiKey"
	AuthMethodOAuth  AuthMethod = "oauth"
)
