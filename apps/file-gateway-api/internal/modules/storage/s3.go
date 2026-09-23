package storage

import (
	"context"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type s3BucketConfig struct {
	Host            string
	Region          string
	AccessKeyID     string
	SecretAccessKey string
	Bucket          string
}

type s3Bucket struct {
	client *s3.Client
	region string
	bucket string
	host   string
}

type s3Provider struct {
	public  s3Bucket
	private s3Bucket
}

func newS3Bucket(config s3BucketConfig) s3Bucket {
	awsCfg := aws.Config{
		Region: config.Region,
		Credentials: aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(
			config.AccessKeyID,
			config.SecretAccessKey,
			"",
		)),
	}

	return s3Bucket{
		client: s3.NewFromConfig(awsCfg),
		region: config.Region,
		bucket: config.Bucket,
		host:   config.Host,
	}
}

func newS3Provider(public s3BucketConfig, private s3BucketConfig) Provider {
	return &s3Provider{
		public:  newS3Bucket(public),
		private: newS3Bucket(private),
	}
}

func (s *s3Provider) Upload(
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

	if err := s3PutObject(ctx, bc, filename, content, contentType, contentLength); err != nil {
		return nil, err
	}

	if !isPublic {
		return &UploadResult{
			Key: filename,
		}, nil
	}

	url := "https://" + bc.bucket + ".s3." + bc.region + ".amazonaws.com/" + filename
	if strings.TrimSpace(bc.host) != "" {
		url = bc.host + "/" + filename
	}

	return &UploadResult{
		URL: url,
		Key: filename,
	}, nil
}

func (s *s3Provider) HealthCheckPublic(ctx context.Context) error {
	return s3HealthCheck(ctx, &s.public)
}

func (s *s3Provider) HealthCheckPrivate(ctx context.Context) error {
	return s3HealthCheck(ctx, &s.private)
}

func s3HealthCheck(ctx context.Context, bc *s3Bucket) error {
	_, err := bc.client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(bc.bucket),
	})
	return err
}

func s3PutObject(
	ctx context.Context,
	bc *s3Bucket,
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
