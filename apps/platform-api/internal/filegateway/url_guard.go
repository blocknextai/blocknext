package filegateway

import (
	"context"
	"net"
	"net/url"
	"strings"
)

func validateDownloadURL(ctx context.Context, rawURL string) error {
	parsed, err := url.Parse(strings.TrimSpace(rawURL))
	if err != nil {
		return ErrInvalidDownloadURL
	}

	switch strings.ToLower(parsed.Scheme) {
	case "http", "https":
	default:
		return ErrInvalidDownloadURL
	}

	host := parsed.Hostname()
	if strings.TrimSpace(host) == "" {
		return ErrInvalidDownloadURL
	}

	if ip := net.ParseIP(host); ip != nil {
		if !isPublicIP(ip) {
			return ErrInvalidDownloadURL
		}
		return nil
	}

	ips, err := net.DefaultResolver.LookupIPAddr(ctx, host)
	if err != nil || len(ips) == 0 {
		return ErrInvalidDownloadURL
	}
	for _, addr := range ips {
		if !isPublicIP(addr.IP) {
			return ErrInvalidDownloadURL
		}
	}

	return nil
}

func isPublicIP(ip net.IP) bool {
	if ip.IsLoopback() ||
		ip.IsPrivate() || // RFC1918 (IPv4) + ULA fc00::/7 (IPv6)
		ip.IsLinkLocalUnicast() || // 169.254.0.0/16 (incl. cloud metadata) + fe80::/10
		ip.IsLinkLocalMulticast() ||
		ip.IsMulticast() ||
		ip.IsUnspecified() {
		return false
	}
	return true
}
