package storage

import (
	"github.com/blocknextai/file-gateway-api/internal/config"
	"github.com/blocknextai/go-packages/apperror"
)

var (
	ErrUnsupportedStorageDriver = apperror.Internal("unsupported storage driver")
)

type Dependencies struct {
	Options config.StorageOptions
}

type Module struct {
	Provider Provider
}

func NewModule(deps Dependencies) (*Module, error) {
	provider, err := buildProvider(deps.Options)
	if err != nil {
		return nil, err
	}
	return &Module{
		Provider: provider,
	}, nil
}

func buildProvider(options config.StorageOptions) (Provider, error) {
	if options.Driver == config.StorageDriverS3 {
		return newS3Provider(
			s3BucketConfig{
				Host:            options.S3.Public.Host,
				Region:          options.S3.Public.Region,
				AccessKeyID:     options.S3.Public.AccessKeyID,
				SecretAccessKey: options.S3.Public.SecretAccessKey,
				Bucket:          options.S3.Public.Bucket,
			},
			s3BucketConfig{
				Host:            options.S3.Private.Host,
				Region:          options.S3.Private.Region,
				AccessKeyID:     options.S3.Private.AccessKeyID,
				SecretAccessKey: options.S3.Private.SecretAccessKey,
				Bucket:          options.S3.Private.Bucket,
			},
		), nil
	}

	if options.Driver == config.StorageDriverLocal {
		return newLocalProvider(
			localPublicBucketConfig{
				Path:    options.Local.Public.Path,
				BaseURL: options.Local.Public.BaseURL,
			},
			localPrivateBucketConfig{
				Path: options.Local.Private.Path,
			},
		), nil
	}

	if options.Driver == config.StorageDriverBunny {
		return newBunnyProvider(
			bunnyBucketConfig{
				StorageHost: options.Bunny.Public.StorageHost,
				StorageZone: options.Bunny.Public.StorageZone,
				AccessKey:   options.Bunny.Public.AccessKey,
				PullZoneURL: options.Bunny.Public.PullZoneURL,
			},
			bunnyBucketConfig{
				StorageHost: options.Bunny.Private.StorageHost,
				StorageZone: options.Bunny.Private.StorageZone,
				AccessKey:   options.Bunny.Private.AccessKey,
				PullZoneURL: options.Bunny.Private.PullZoneURL,
			},
		), nil
	}

	return nil, ErrUnsupportedStorageDriver
}
