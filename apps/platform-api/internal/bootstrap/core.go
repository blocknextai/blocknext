package bootstrap

import (
	"context"
	"database/sql"
	"log/slog"

	"github.com/blocknextai/go-packages/cache"
	"github.com/blocknextai/go-packages/database"
	"github.com/blocknextai/go-packages/secretmanager"
	cacheInfrastructure "github.com/blocknextai/platform-api/internal/cache/infrastructure"
	"github.com/blocknextai/platform-api/internal/config"
	"github.com/blocknextai/platform-api/internal/eventbus"
	"github.com/blocknextai/platform-api/internal/filegateway"
)

type Core struct {
	Config *config.SharedConfig

	DB                 *sql.DB
	CacheService       cache.Service
	TransactionManager database.TransactionManager
	SecretManager      secretmanager.SecretManager
	FileGateway        filegateway.FileGateway
	EventBus           *eventbus.Module
}

type Option func(*options)

type options struct {
	dbMaxOpenConns int
	dbMaxIdleConns int
}

func WithDBPool(maxOpen, maxIdle int) Option {
	return func(o *options) {
		o.dbMaxOpenConns = maxOpen
		o.dbMaxIdleConns = maxIdle
	}
}

func NewCore(cfg *config.SharedConfig, opts ...Option) (*Core, error) {
	o := &options{
		dbMaxOpenConns: cfg.Database.MaxOpenConns,
		dbMaxIdleConns: cfg.Database.MaxIdleConns,
	}
	for _, apply := range opts {
		apply(o)
	}

	db, err := database.NewDB(
		cfg.Database.DSN(),
		o.dbMaxOpenConns,
		o.dbMaxIdleConns,
		cfg.Database.ConnMaxLifetime,
		cfg.Database.ConnMaxIdleTime,
	)
	if err != nil {
		return nil, err
	}

	cacheService, err := cacheInfrastructure.NewCacheService(cfg.Cache)
	if err != nil {
		return nil, err
	}

	secretManager, err := secretmanager.New(cfg.SecretManager.SecretKey)
	if err != nil {
		return nil, err
	}

	fileGateway := filegateway.NewFileGateway(
		cfg.FileGateway.BaseURL,
		cfg.FileGateway.Auth.ServiceKey,
	)

	transactionManager := database.NewTransactionManager(db, cfg.Database.QueryTimeout)

	eventBusModule := eventbus.NewModule(eventbus.Dependencies{
		DB:                 db,
		TransactionManager: transactionManager,

		Options: cfg.EventBus,
	})

	return &Core{
		Config:             cfg,
		DB:                 db,
		CacheService:       cacheService,
		TransactionManager: transactionManager,
		SecretManager:      secretManager,
		FileGateway:        fileGateway,
		EventBus:           eventBusModule,
	}, nil
}

func (c *Core) Health(ctx context.Context) error {
	if err := c.DB.PingContext(ctx); err != nil {
		return err
	}
	return c.CacheService.Ping(ctx)
}

func (c *Core) Shutdown() error {
	var firstErr error
	if err := c.CacheService.Close(); err != nil {
		slog.Error("failed to close cache service", "error", err)
		firstErr = err
	}
	if err := c.DB.Close(); err != nil {
		slog.Error("failed to close database", "error", err)
		if firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}
