package local

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"

	storageDomain "github.com/blocknextai/file-gateway-api/internal/storage/domain/storage"
	"github.com/blocknextai/go-packages/apperror"
)

var (
	ErrInvalidFilename  = apperror.Validation("invalid filename")
	ErrPathNotDirectory = apperror.Internal("path is not a directory")
)

type PublicBucketConfig struct {
	Path    string
	BaseURL string
}

type PrivateBucketConfig struct {
	Path string
}

type bucketClient struct {
	path    string
	baseURL string
}

type storageProvider struct {
	public  bucketClient
	private bucketClient
}

func NewStorageProvider(public PublicBucketConfig, private PrivateBucketConfig) storageDomain.Provider {
	_ = os.MkdirAll(public.Path, 0o755)
	_ = os.MkdirAll(private.Path, 0o755)

	return &storageProvider{
		public: bucketClient{
			path:    public.Path,
			baseURL: strings.TrimRight(public.BaseURL, "/"),
		},
		private: bucketClient{
			path: private.Path,
		},
	}
}

func (s *storageProvider) Upload(
	ctx context.Context,
	isPublic bool,
	filename string,
	content io.Reader,
	contentType string,
	contentLength int64,
) (*storageDomain.UploadResult, error) {
	bc := &s.private
	if isPublic {
		bc = &s.public
	}

	if err := writeFile(ctx, bc, filename, content); err != nil {
		return nil, err
	}

	if isPublic {
		return &storageDomain.UploadResult{
			URL: bc.baseURL + "/" + filename,
			Key: filename,
		}, nil
	}
	return &storageDomain.UploadResult{
		Key: filename,
	}, nil
}

func (s *storageProvider) HealthCheckPublic(ctx context.Context) error {
	return healthCheck(&s.public)
}

func (s *storageProvider) HealthCheckPrivate(ctx context.Context) error {
	return healthCheck(&s.private)
}

func writeFile(ctx context.Context, bc *bucketClient, key string, content io.Reader) (err error) {
	if err := ctx.Err(); err != nil {
		return err
	}

	target, err := safeJoin(bc.path, key)
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(target), 0o755); err != nil {
		return err
	}

	f, err := os.Create(target)
	if err != nil {
		return err
	}
	defer func() {
		if cerr := f.Close(); cerr != nil && err == nil {
			err = cerr
		}
		if err != nil {
			_ = os.Remove(target)
		}
	}()

	_, err = io.Copy(f, content)
	return err
}

func healthCheck(bc *bucketClient) (err error) {
	info, err := os.Stat(bc.path)
	if err != nil {
		return err
	}
	if !info.IsDir() {
		return ErrPathNotDirectory
	}

	f, err := os.CreateTemp(bc.path, ".healthcheck-*")
	if err != nil {
		return err
	}
	name := f.Name()
	defer func() {
		if cerr := f.Close(); cerr != nil && err == nil {
			err = cerr
		}
		if rerr := os.Remove(name); rerr != nil && err == nil {
			err = rerr
		}
	}()
	return nil
}

func safeJoin(base, name string) (string, error) {
	cleaned := filepath.Join(base, name)
	rel, err := filepath.Rel(base, cleaned)
	if err != nil {
		return "", err
	}
	if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", ErrInvalidFilename
	}
	return cleaned, nil
}
