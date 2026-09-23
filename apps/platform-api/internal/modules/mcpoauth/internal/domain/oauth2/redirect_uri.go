package oauth2

import (
	"net"
	"net/url"
	"strings"
)

func IsValidRedirectURI(rawURI string) bool {
	trimmed := strings.TrimSpace(rawURI)
	if trimmed == "" || strings.Contains(trimmed, "#") {
		return false
	}

	parsed, err := url.Parse(trimmed)
	if err != nil {
		return false
	}

	switch strings.ToLower(parsed.Scheme) {
	case "https":
		return strings.TrimSpace(parsed.Host) != ""
	case "http":
		return isLoopbackHost(parsed.Hostname())
	default:
		return strings.TrimSpace(parsed.Scheme) != "" &&
			(strings.TrimSpace(parsed.Opaque) != "" || strings.TrimSpace(parsed.Host) != "" || strings.TrimSpace(parsed.Path) != "")
	}
}

func BuildRedirectURI(redirectURI string, params map[string]string) (string, error) {
	parsed, err := url.Parse(redirectURI)
	if err != nil {
		return "", ErrMalformedRedirectURI.WithCause(err)
	}

	query := parsed.Query()
	for key, value := range params {
		if strings.TrimSpace(value) == "" {
			continue
		}
		query.Set(key, value)
	}
	parsed.RawQuery = query.Encode()

	return parsed.String(), nil
}

func isLoopbackHost(host string) bool {
	if host == "localhost" {
		return true
	}

	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}
