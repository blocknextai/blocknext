# internal/common

Shared domain primitives reused across bounded contexts. Intentionally minimal —
only put something here when it is genuinely cross-context.

## Scope

- **`LimitedReader`** — an `io.Reader` wrapper that caps how many bytes may be read
  and returns a caller-supplied error once the limit is exceeded (used by the
  download context to enforce max response size).
- **Shared errors** — generic request errors (e.g. `ErrInvalidRequest`) not owned by
  a single context.

## Layout

```
domain/reader.go     # LimitedReader
domain/errors.go     # Shared domain errors
```
