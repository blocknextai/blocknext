# internal/auth

Authentication bounded context. Owns token issuance and the request-guarding
middleware that protects the upload and download routes.

## Scope

- **Token issuance** — mints short-lived JWTs bound to a random session ID via
  `POST /auth/token`. This endpoint is intentionally **public** so trusted UI
  clients can obtain a token; it is only IP rate-limited.
- **Request authentication** — middleware that accepts either an
  `X-Service-Key: <key>` header (server-to-server) or an
  `Authorization: Bearer <jwt>` token.
- **Service-key matching** — `MatchServiceKey` lets the global rate limiter skip
  trusted service-key callers.

## Layout

```
domain/token/        # Service interface (token contract) + domain errors
infrastructure/token/ # JWT-backed Service implementation (go-packages/auth/jwt)
presentation/        # /auth/token route, auth middleware, request/response types
module.go            # Wires the JWT service + middleware; exposes Register()
```

## Public surface

`auth.NewModule(Dependencies{ServiceKey, JWTService})` returns a `Module` exposing
`TokenService`, the `Middleware` handler, and the `MatchServiceKey` predicate.
`Module.Register(router)` mounts the public token route.
