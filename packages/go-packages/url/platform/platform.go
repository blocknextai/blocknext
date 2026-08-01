// Package platform detects the social or web platform a URL belongs to from its
// domain.
package platform

import (
	"net/url"
	"strings"
)

const defaultPlatform = "website"

var platformMapping = map[string]string{
	"twitter":   "x",
	"x":         "x",
	"linkedin":  "linkedin",
	"github":    "github",
	"instagram": "instagram",
	"facebook":  "facebook",
	"youtube":   "youtube",
	"tiktok":    "tiktok",
	"discord":   "discord",
	"telegram":  "telegram",
	"whatsapp":  "whatsapp",
	"snapchat":  "snapchat",
	"pinterest": "pinterest",
	"reddit":    "reddit",
	"twitch":    "twitch",
	"medium":    "medium",
	"behance":   "behance",
	"dribbble":  "dribbble",
	"spotify":   "spotify",
	"slack":     "slack",
	"threads":   "threads",
}

// Detect returns the platform identifier for the given URL based on its domain,
// or the default platform when the URL cannot be parsed or is unrecognized.
func Detect(rawurl string) string {
	parsed, err := url.Parse(rawurl)
	if err != nil {
		return defaultPlatform
	}

	hostname := strings.ToLower(parsed.Hostname())
	parts := strings.Split(hostname, ".")
	if len(parts) < 2 {
		return defaultPlatform
	}

	domain := parts[len(parts)-2]
	if p, exists := platformMapping[domain]; exists {
		return p
	}

	return defaultPlatform
}
