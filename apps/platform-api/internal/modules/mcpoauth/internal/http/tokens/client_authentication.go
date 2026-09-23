package tokens

import (
	"net/url"
	"strings"

	"github.com/blocknextai/go-packages/base64"
)

const (
	basicAuthenticationPrefix = "Basic "
)

func clientCredentials(authorization string, clientID string, clientSecret string) (string, string) {
	encoded, ok := strings.CutPrefix(authorization, basicAuthenticationPrefix)
	if !ok {
		return clientID, clientSecret
	}

	decoded, err := base64.Decode(strings.TrimSpace(encoded))
	if err != nil {
		return clientID, clientSecret
	}

	basicClientID, basicClientSecret, found := strings.Cut(string(decoded), ":")
	if !found {
		return clientID, clientSecret
	}

	return formDecode(basicClientID), formDecode(basicClientSecret)
}

func formDecode(value string) string {
	decoded, err := url.QueryUnescape(value)
	if err != nil {
		return value
	}
	return decoded
}
