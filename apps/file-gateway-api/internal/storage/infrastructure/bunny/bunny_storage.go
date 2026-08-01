package bunny

import (
	"context"
	"io"
	"strings"
	"time"

	storageDomain "github.com/blocknextai/file-gateway-api/internal/storage/domain/storage"
	"github.com/blocknextai/go-packages/apperror"
	"github.com/blocknextai/go-packages/httpclient"
)

var (
	ErrUploadFailed      = apperror.Internal("bunny storage upload failed")
	ErrHealthCheckFailed = apperror.Internal("bunny storage health check failed")
)

type BucketConfig struct {
	StorageHost string
	StorageZone string
	AccessKey   string
	PullZoneURL string
}

type bucketClient struct {
	client      *httpclient.Client
	storageZone string
	pullZoneURL string
}

type storageProvider struct {
	public  bucketClient
	private bucketClient
}

func newBucketClient(config BucketConfig) bucketClient {
	client := httpclient.NewClientBuilder().
		BaseURL(strings.TrimRight(config.StorageHost, "/")).
		Header("AccessKey", config.AccessKey).
		Timeout(5*time.Minute).
		RetryConfig(0, 0).
		Build()

	return bucketClient{
		client:      client,
		storageZone: config.StorageZone,
		pullZoneURL: strings.TrimRight(config.PullZoneURL, "/"),
	}
}

func NewStorageProvider(public BucketConfig, private BucketConfig) storageDomain.Provider {
	return &storageProvider{
		public:  newBucketClient(public),
		private: newBucketClient(private),
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

	if err := putObject(ctx, bc, filename, content, contentType, contentLength); err != nil {
		return nil, err
	}

	if isPublic {
		return &storageDomain.UploadResult{
			URL: bc.pullZoneURL + "/" + filename,
			Key: filename,
		}, nil
	}
	return &storageDomain.UploadResult{
		Key: filename,
	}, nil
}

func (s *storageProvider) HealthCheckPublic(ctx context.Context) error {
	return healthCheck(ctx, &s.public)
}

func (s *storageProvider) HealthCheckPrivate(ctx context.Context) error {
	return healthCheck(ctx, &s.private)
}

func putObject(ctx context.Context, bc *bucketClient, key string, content io.Reader, contentType string, contentLength int64) error {
	req := bc.client.Put("/" + bc.storageZone + "/" + key).
		Context(ctx).
		BodyReader(content).
		ContentLength(contentLength)

	if strings.TrimSpace(contentType) != "" {
		req = req.ContentType(contentType)
	}

	resp, err := req.DoRaw()
	if err != nil {
		return err
	}
	if !resp.IsSuccess() {
		return ErrUploadFailed
	}

	return nil
}

func healthCheck(ctx context.Context, bc *bucketClient) error {
	resp, err := bc.client.Get("/" + bc.storageZone + "/").Context(ctx).DoRaw()
	if err != nil {
		return err
	}
	if !resp.IsSuccess() {
		return ErrHealthCheckFailed
	}
	return nil
}
