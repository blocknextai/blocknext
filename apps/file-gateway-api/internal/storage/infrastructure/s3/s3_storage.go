package s3

import (
	"context"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	storageDomain "github.com/blocknextai/file-gateway-api/internal/storage/domain/storage"
)

type BucketConfig struct {
	Host            string
	Region          string
	AccessKeyID     string
	SecretAccessKey string
	Bucket          string
}

type bucketClient struct {
	client *s3.Client
	region string
	bucket string
	host   string
}

type storageProvider struct {
	public  bucketClient
	private bucketClient
}

func newBucketClient(config BucketConfig) bucketClient {
	awsCfg := aws.Config{
		Region: config.Region,
		Credentials: aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(
			config.AccessKeyID,
			config.SecretAccessKey,
			"",
		)),
	}

	return bucketClient{
		client: s3.NewFromConfig(awsCfg),
		region: config.Region,
		bucket: config.Bucket,
		host:   config.Host,
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

	if !isPublic {
		return &storageDomain.UploadResult{
			Key: filename,
		}, nil
	}

	url := "https://" + bc.bucket + ".s3." + bc.region + ".amazonaws.com/" + filename
	if strings.TrimSpace(bc.host) != "" {
		url = bc.host + "/" + filename
	}

	return &storageDomain.UploadResult{
		URL: url,
		Key: filename,
	}, nil
}

func (s *storageProvider) HealthCheckPublic(ctx context.Context) error {
	return healthCheck(ctx, &s.public)
}

func (s *storageProvider) HealthCheckPrivate(ctx context.Context) error {
	return healthCheck(ctx, &s.private)
}

func healthCheck(ctx context.Context, bc *bucketClient) error {
	_, err := bc.client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(bc.bucket),
	})
	return err
}

func putObject(
	ctx context.Context,
	bc *bucketClient,
	key string,
	content io.Reader,
	contentType string,
	contentLength int64,
) error {
	_, err := bc.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(bc.bucket),
		Key:           aws.String(key),
		Body:          content,
		ContentType:   aws.String(contentType),
		ContentLength: aws.Int64(contentLength),
	})
	return err
}
