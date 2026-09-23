package oauth2

import (
	"net/url"
	"strings"
)

func CanonicalizeResource(rawResource string) string {
	trimmed := strings.TrimSpace(rawResource)
	if trimmed == "" || strings.Contains(trimmed, "#") {
		return ""
	}

	parsed, err := url.Parse(trimmed)
	if err != nil {
		return ""
	}

	scheme := strings.ToLower(parsed.Scheme)
	if scheme != "http" && scheme != "https" {
		return ""
	}
	if strings.TrimSpace(parsed.Host) == "" {
		return ""
	}

	parsed.Scheme = scheme
	parsed.Host = strings.ToLower(parsed.Host)
	parsed.Path = strings.TrimSuffix(parsed.Path, "/")

	return parsed.String()
}

func ResolveResource(resourceBaseURL string, rawResource string) (string, bool) {
	base := CanonicalizeResource(resourceBaseURL)
	if strings.TrimSpace(base) == "" {
		return "", false
	}

	if strings.TrimSpace(rawResource) == "" {
		return base, true
	}

	resource := CanonicalizeResource(rawResource)
	if resource == base || strings.HasPrefix(resource, base+"/") {
		return resource, true
	}

	return "", false
}
