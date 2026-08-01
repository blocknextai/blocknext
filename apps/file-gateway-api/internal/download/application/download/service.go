package download

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/netip"
	"net/url"
	"slices"
	"strconv"
	"strings"
	"time"

	commonDomain "github.com/blocknextai/file-gateway-api/internal/common/domain"
	downloadDomain "github.com/blocknextai/file-gateway-api/internal/download/domain"
	downloadDomainDownload "github.com/blocknextai/file-gateway-api/internal/download/domain/download"
	"github.com/blocknextai/go-packages/file"
	"github.com/blocknextai/go-packages/httpclient"
)

var (
	blockedPrefixes = []netip.Prefix{
		netip.MustParsePrefix("100.64.0.0/10"),      // RFC 6598 carrier-grade NAT / shared address space
		netip.MustParsePrefix("192.0.0.0/24"),       // RFC 6890 IETF protocol assignments
		netip.MustParsePrefix("192.0.2.0/24"),       // TEST-NET-1
		netip.MustParsePrefix("198.18.0.0/15"),      // RFC 2544 benchmarking
		netip.MustParsePrefix("198.51.100.0/24"),    // TEST-NET-2
		netip.MustParsePrefix("203.0.113.0/24"),     // TEST-NET-3
		netip.MustParsePrefix("240.0.0.0/4"),        // reserved
		netip.MustParsePrefix("255.255.255.255/32"), // limited broadcast
		netip.MustParsePrefix("64:ff9b::/96"),       // NAT64 well-known prefix
		netip.MustParsePrefix("2002::/16"),          // 6to4
		netip.MustParsePrefix("2001::/32"),          // Teredo
		netip.MustParsePrefix("fec0::/10"),          // deprecated IPv6 site-local
	}
)

type service struct {
	client  *httpclient.Client
	maxSize int64
}

func NewService(maxSize int64, timeout time.Duration) downloadDomainDownload.Service {
	return &service{
		client: httpclient.NewClientBuilder().
			Timeout(timeout).
			RetryConfig(0, 0).
			NoRedirect().
			Transport(safeTransport(timeout)).
			Build(),
		maxSize: maxSize,
	}
}

func safeTransport(timeout time.Duration) http.RoundTripper {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	dialer := &net.Dialer{
		Timeout:   timeout,
		KeepAlive: 30 * time.Second,
	}
	transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		host, port, err := net.SplitHostPort(addr)
		if err != nil {
			return nil, err
		}

		ips, err := net.DefaultResolver.LookupIP(ctx, "ip", host)
		if err != nil {
			return nil, err
		}

		if slices.ContainsFunc(ips, isBlockedIP) {
			return nil, downloadDomain.ErrBlockedIP
		}

		if len(ips) == 0 {
			return nil, downloadDomain.ErrInvalidURL
		}

		return dialer.DialContext(ctx, network, net.JoinHostPort(ips[0].String(), port))
	}
	return transport
}

func (s *service) Download(ctx context.Context, rawURL string) (*downloadDomainDownload.FileInfo, error) {
	parsedURL, err := url.Parse(rawURL)
	if err != nil {
		return nil, downloadDomain.ErrInvalidURL
	}

	if err := validateURL(parsedURL); err != nil {
		return nil, err
	}

	resp, err := s.client.Get(rawURL).Context(ctx).DoStream()
	if err != nil {
		if strings.Contains(err.Error(), "timeout") {
			return nil, downloadDomain.ErrDownloadTimeout
		}
		slog.Error("download request failed", "url", rawURL, "error", err)
		return nil, downloadDomain.ErrDownloadFailed
	}

	if !resp.IsSuccess() {
		closeBody(resp.BodyReader)
		return nil, downloadDomain.ErrDownloadFailed
	}

	contentLength := parseContentLength(resp.Headers.Get("Content-Length"))
	if contentLength > s.maxSize {
		closeBody(resp.BodyReader)
		return nil, downloadDomain.ErrFileTooLarge
	}

	return s.buildFileInfo(resp.BodyReader, contentLength, rawURL)
}

func (s *service) buildFileInfo(body io.ReadCloser, contentLength int64, rawURL string) (*downloadDomainDownload.FileInfo, error) {
	limited := commonDomain.NewLimitedReader(body, s.maxSize, downloadDomain.ErrFileTooLarge)

	header := make([]byte, 512)
	n, err := io.ReadFull(limited, header)
	if err != nil && err != io.EOF && err != io.ErrUnexpectedEOF {
		closeBody(body)
		if errors.Is(err, downloadDomain.ErrFileTooLarge) {
			return nil, downloadDomain.ErrFileTooLarge
		}
		return nil, downloadDomain.ErrDownloadFailed
	}
	header = header[:n]

	mime := file.MIMEType(header)
	combined := io.MultiReader(bytes.NewReader(header), limited)

	return &downloadDomainDownload.FileInfo{
		Filename:    file.ExtractFilename(rawURL),
		Type:        file.Category(mime.Extension()),
		Ext:         mime.Extension(),
		ContentType: mime.String(),
		Size:        contentLength,
		BodyReader:  &closerReader{reader: combined, closer: body},
	}, nil
}

func parseContentLength(v string) int64 {
	if strings.TrimSpace(v) == "" {
		return -1
	}
	n, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return -1
	}
	return n
}

func validateURL(u *url.URL) error {
	if u.Scheme != "https" {
		return downloadDomain.ErrInvalidScheme
	}

	ips, err := net.LookupIP(u.Hostname())
	if err != nil {
		return downloadDomain.ErrInvalidURL
	}

	if slices.ContainsFunc(ips, isBlockedIP) {
		return downloadDomain.ErrBlockedIP
	}
	return nil
}

func isBlockedIP(ip net.IP) bool {
	if ip.IsLoopback() || ip.IsPrivate() || ip.IsUnspecified() {
		return true
	}
	if ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
		return true
	}

	addr, ok := netip.AddrFromSlice(ip)
	if !ok {
		return true
	}
	addr = addr.Unmap()
	return slices.ContainsFunc(blockedPrefixes, func(p netip.Prefix) bool {
		return p.Contains(addr)
	})
}

func closeBody(body io.Closer) {
	if err := body.Close(); err != nil {
		slog.Error("failed to close response body", "error", err)
	}
}

type closerReader struct {
	reader io.Reader
	closer io.Closer
}

func (c *closerReader) Read(p []byte) (int, error) {
	return c.reader.Read(p)
}

func (c *closerReader) Close() error {
	return c.closer.Close()
}
