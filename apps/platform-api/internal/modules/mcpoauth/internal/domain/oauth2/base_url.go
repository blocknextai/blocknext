package oauth2

import (
	"net/url"
	"strings"
)

func IsValidBaseURL(rawURL string) bool {
	parsed, err := url.Parse(strings.TrimSpace(rawURL))
	if err != nil {
		return false
	}

	return parsed.Scheme != "" &&
		parsed.Host != "" &&
		strings.Trim(parsed.Path, "/") == "" &&
		parsed.RawQuery == "" &&
		parsed.Fragment == ""
}
