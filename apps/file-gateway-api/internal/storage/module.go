package storage

import (
	"github.com/blocknextai/file-gateway-api/internal/config"
	storageDomain "github.com/blocknextai/file-gateway-api/internal/storage/domain/storage"
	"github.com/blocknextai/file-gateway-api/internal/storage/infrastructure/bunny"
	"github.com/blocknextai/file-gateway-api/internal/storage/infrastructure/local"
	"github.com/blocknextai/file-gateway-api/internal/storage/infrastructure/s3"
	"github.com/blocknextai/go-packages/apperror"
)

var (
	ErrUnsupportedStorageDriver = apperror.Internal("unsupported storage driver")
)

type Dependencies struct {
	Options config.StorageOptions
}

type Module struct {
	Provider storageDomain.Provider
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

func buildProvider(options config.StorageOptions) (storageDomain.Provider, error) {
	if options.Driver == config.StorageDriverS3 {
		return s3.NewStorageProvider(
			s3.BucketConfig{
				Host:            options.S3.Public.Host,
				Region:          options.S3.Public.Region,
				AccessKeyID:     options.S3.Public.AccessKeyID,
				SecretAccessKey: options.S3.Public.SecretAccessKey,
				Bucket:          options.S3.Public.Bucket,
			},
			s3.BucketConfig{
				Host:            options.S3.Private.Host,
				Region:          options.S3.Private.Region,
				AccessKeyID:     options.S3.Private.AccessKeyID,
				SecretAccessKey: options.S3.Private.SecretAccessKey,
				Bucket:          options.S3.Private.Bucket,
			},
		), nil
	}

	if options.Driver == config.StorageDriverLocal {
		return local.NewStorageProvider(
			local.PublicBucketConfig{
				Path:    options.Local.Public.Path,
				BaseURL: options.Local.Public.BaseURL,
			},
			local.PrivateBucketConfig{
				Path: options.Local.Private.Path,
			},
		), nil
	}

	if options.Driver == config.StorageDriverBunny {
		return bunny.NewStorageProvider(
			bunny.BucketConfig{
				StorageHost: options.Bunny.Public.StorageHost,
				StorageZone: options.Bunny.Public.StorageZone,
				AccessKey:   options.Bunny.Public.AccessKey,
				PullZoneURL: options.Bunny.Public.PullZoneURL,
			},
			bunny.BucketConfig{
				StorageHost: options.Bunny.Private.StorageHost,
				StorageZone: options.Bunny.Private.StorageZone,
				AccessKey:   options.Bunny.Private.AccessKey,
				PullZoneURL: options.Bunny.Private.PullZoneURL,
			},
		), nil
	}

	return nil, ErrUnsupportedStorageDriver
}
