package filegateway

import (
	"context"
	"io"
)

type DownloadResult struct {
	Filename    string
	Type        string
	Ext         string
	ContentType string
	Size        int64
	DataReader  io.ReadCloser
}

type UploadResult struct {
	URL string
}

type FileGateway interface {
	DownloadFile(ctx context.Context, url string) (*DownloadResult, error)
	UploadFile(ctx context.Context, uploadID string, fileName string, reader io.Reader) (*UploadResult, error)
}
