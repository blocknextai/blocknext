# internal/cache

Cache infrastructure. Provides the shared `cache.Service` used for rate-limiter
storage and readiness pings.

## Scope

- **Service factory** — `NewCacheService` builds a `cache.Service` from `CACHE_*`
  configuration. Currently supports the `redis` type (with pool tuning); an unknown
  type returns `ErrInvalidCacheType`.

## Layout

```
infrastructure/factory.go  # NewCacheService — driver selection + Redis pool config
```

This package has no domain layer; it is pure wiring around
`go-packages/cache`. The cache instance is created in `main` and shared with the
rate-limiter middleware.
