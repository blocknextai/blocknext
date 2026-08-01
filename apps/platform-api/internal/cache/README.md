# Cache

> Factory that builds the shared cache service from config. Infrastructure, not a business bounded context.

## Responsibility
Selects and constructs a `cache.Service` (from `go-packages/cache`) based on `config.CacheOptions`. The cache abstraction itself lives in `go-packages`; this package only owns the config-driven wiring that picks an implementation.

## What it provides
- `NewCacheService(config.CacheOptions) (cache.Service, error)` — returns a Redis-backed `cache.Service` when `Type == CacheTypeRedis`, wiring address/password/DB and the connection-pool options.
- `ErrInvalidCacheType` — returned for any unsupported cache type.

## Used by
`bootstrap.NewCore`, which stores the result as `Core.CacheService` and shares it across modules (account, organizations, oauth, taskrunner, …) and the health check (`CacheService.Ping`).

## Notes
- Only Redis is implemented today; the factory is structured to allow other backends without touching consumers.
- The returned `cache.Service` is the go-packages interface — this package adds no methods of its own.

## Layout
- `infrastructure/factory.go` — the `NewCacheService` factory and `ErrInvalidCacheType`.
