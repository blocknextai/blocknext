package storage

import (
	"context"
	"io"
)

type UploadResult struct {
	URL string `json:"url,omitempty"`
	Key string `json:"key,omitempty"`
}

type Provider interface {
	Upload(ctx context.Context, isPublic bool, filename string, content io.Reader, contentType string, contentLength int64) (*UploadResult, error)
	HealthCheckPublic(ctx context.Context) error
	HealthCheckPrivate(ctx context.Context) error
}
