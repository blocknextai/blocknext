package config

type StorageDriver string

const (
	StorageDriverLocal StorageDriver = "local"
	StorageDriverS3    StorageDriver = "s3"
	StorageDriverBunny StorageDriver = "bunny"
)

type StorageOptions struct {
	Driver StorageDriver `env:"DRIVER"`

	Local LocalOptions `envPrefix:"LOCAL_"`
	S3    S3Options    `envPrefix:"S3_"`
	Bunny BunnyOptions `envPrefix:"BUNNY_"`
}

type LocalOptions struct {
	Public  LocalPublicOptions  `envPrefix:"PUBLIC_"`
	Private LocalPrivateOptions `envPrefix:"PRIVATE_"`
}

type LocalPublicOptions struct {
	Path    string `env:"PATH"`
	BaseURL string `env:"BASE_URL"`
}

type LocalPrivateOptions struct {
	Path string `env:"PATH"`
}

type S3Options struct {
	Public  S3BucketOptions `envPrefix:"PUBLIC_"`
	Private S3BucketOptions `envPrefix:"PRIVATE_"`
}

type S3BucketOptions struct {
	Host            string `env:"HOST"`
	AccessKeyID     string `env:"ACCESS_KEY_ID"`
	SecretAccessKey string `env:"SECRET_ACCESS_KEY"`
	Region          string `env:"REGION"`
	Bucket          string `env:"BUCKET"`
}

type BunnyOptions struct {
	Public  BunnyZoneOptions `envPrefix:"PUBLIC_"`
	Private BunnyZoneOptions `envPrefix:"PRIVATE_"`
}

type BunnyZoneOptions struct {
	StorageHost string `env:"STORAGE_HOST"`
	StorageZone string `env:"STORAGE_ZONE"`
	AccessKey   string `env:"ACCESS_KEY"`
	PullZoneURL string `env:"PULL_ZONE_URL"`
}
