package storage

import (
	"context"
	"io"
	"strings"
	"time"

	"github.com/blocknextai/go-packages/apperror"
	"github.com/blocknextai/go-packages/httpclient"
)

var (
	ErrUploadFailed      = apperror.Internal("bunny storage upload failed")
	ErrHealthCheckFailed = apperror.Internal("bunny storage health check failed")
)

type bunnyBucketConfig struct {
	StorageHost string
	StorageZone string
	AccessKey   string
	PullZoneURL string
}

type bunnyBucket struct {
	client      *httpclient.Client
	storageZone string
	pullZoneURL string
}

type bunnyProvider struct {
	public  bunnyBucket
	private bunnyBucket
}

func newBunnyBucket(config bunnyBucketConfig) bunnyBucket {
	client := httpclient.NewClientBuilder().
		BaseURL(strings.TrimRight(config.StorageHost, "/")).
		Header("AccessKey", config.AccessKey).
		Timeout(5*time.Minute).
		RetryConfig(0, 0).
		Build()

	return bunnyBucket{
		client:      client,
		storageZone: config.StorageZone,
		pullZoneURL: strings.TrimRight(config.PullZoneURL, "/"),
	}
}

func newBunnyProvider(public bunnyBucketConfig, private bunnyBucketConfig) Provider {
	return &bunnyProvider{
		public:  newBunnyBucket(public),
		private: newBunnyBucket(private),
	}
}

func (s *bunnyProvider) Upload(
	ctx context.Context,
	isPublic bool,
	filename string,
	content io.Reader,
	contentType string,
	contentLength int64,
) (*UploadResult, error) {
	bc := &s.private
	if isPublic {
		bc = &s.public
	}

	if err := bunnyPutObject(ctx, bc, filename, content, contentType, contentLength); err != nil {
		return nil, err
	}

	if isPublic {
		return &UploadResult{
			URL: bc.pullZoneURL + "/" + filename,
			Key: filename,
		}, nil
	}
	return &UploadResult{
		Key: filename,
	}, nil
}

func (s *bunnyProvider) HealthCheckPublic(ctx context.Context) error {
	return bunnyHealthCheck(ctx, &s.public)
}

func (s *bunnyProvider) HealthCheckPrivate(ctx context.Context) error {
	return bunnyHealthCheck(ctx, &s.private)
}

func bunnyPutObject(ctx context.Context, bc *bunnyBucket, key string, content io.Reader, contentType string, contentLength int64) error {
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

func bunnyHealthCheck(ctx context.Context, bc *bunnyBucket) error {
	resp, err := bc.client.Get("/" + bc.storageZone + "/").Context(ctx).DoRaw()
	if err != nil {
		return err
	}
	if !resp.IsSuccess() {
		return ErrHealthCheckFailed
	}
	return nil
}
