package usecases

import (
	"context"
	"log/slog"
	"net/url"
	"strings"

	downloadDomain "github.com/blocknextai/file-gateway-api/internal/modules/download/internal/domain"
)

func (s *Service) Download(ctx context.Context, rawURL string) (*downloadDomain.FileInfo, error) {
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
