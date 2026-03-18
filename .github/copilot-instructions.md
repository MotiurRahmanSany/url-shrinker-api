# URL Shrinker API — Copilot Instructions

## Project Overview

Go REST API for shortening URLs with JWT auth, click analytics, and Redis caching.
Module: `github.com/MotiurRahmanSany/url-shrinker-api`

## Commands

```bash
make run           # docker compose up + air (live reload)
make stop          # stop containers
make migrate-up    # run goose migrations
make migrate-down  # rollback last migration
make sqlc          # regenerate db code from sql/queries/
make flush-redis   # clear all Redis data
make seed-admin    # run cmd/seed (stub, not yet implemented)
```

Air live-reload watches `.go/.tpl/.html` and rebuilds to `./tmp/main`.
The app **panics on startup if `.env` is missing** (godotenv.Load is non-optional).

## Architecture

```
HTTP Handlers  →  Services  →  Repositories  →  sqlc Queries (internal/db/)
                     ↓
                Redis Cache
```

**Dependency wiring** is all in `cmd/api/serve.go`:
`pgxpool` → `db.New()` → repositories → services → handlers → router

**Middleware chain:** `Logger` wraps everything; `AuthMiddleware` wraps protected routes; `AdminOnly` (defined, currently unused in routes).

## Layer Conventions

### Handlers (`internal/api/handlers/`)
- Decode request body with `json.NewDecoder(r.Body).Decode(&req)`.
- Extract auth context: `r.Context().Value(middleware.UserContextKey).(string)`.
- Return responses exclusively via `response.Success()` / `response.Error()`.
- Map service errors with `errors.Is()` to appropriate HTTP status codes.

### Services (`internal/service/`)
- Define named sentinel errors for business rule failures (e.g. `ErrEmailAlreadyInUse`).
- No HTTP concerns — return domain types and errors only.

### Repositories (`internal/repository/`)
- Accept and return **domain types** (`internal/domain/`), never sqlc/pgtype types.
- Always convert pgtype → Go types before returning:
  - `pgtype.UUID` → `row.ID.String()`
  - `pgtype.Timestamp` / `pgtype.Timestamptz` → `row.CreatedAt.Time`
  - `pgtype.Int4` → `row.MaxClicks.Int32` (check `.Valid` for nullable)
- Converting Go → pgtype for query params: use `.Scan()`:
  ```go
  var pgID pgtype.UUID
  pgID.Scan(userID) // string → pgtype.UUID
  ```

### Domain Types (`internal/domain/`)
- Plain Go structs with `time.Time`, `string` IDs, pointer types for nullable fields (`*int32`, `*time.Time`).
- These are what handlers and services work with — never expose sqlc types above the repository layer.

## Response Envelope

All responses use this structure:
```go
response.Success(w, http.StatusOK, "message", data)
response.Error(w, http.StatusBadRequest, "message", nil)
```
JSON shape: `{ success, message, data?, error? }`

For paginated lists use `response.NewPaginatedData(items, page, limit, total)`.

## Auth Context Keys

Defined in `internal/api/middleware/auth_middleware.go`:
- `middleware.UserContextKey` → `string` (UUID of authenticated user)
- `middleware.RoleContextKey` → `string` (`"admin"` or `"student"`)

## Database

- **pgxpool** for connection pooling; pool is created in `internal/database/connection.go`.
- **sqlc** generates `internal/db/` from `sql/queries/`. After changing any `.sql` query file, run `make sqlc`.
- sqlc config: `sqlc.yaml` (engine: postgresql, emits JSON tags + interface).

## Redis Cache

Interface in `internal/cache/redis.go`:
```go
Get(ctx, key) (string, error)       // returns ErrCacheMiss on miss
Set(ctx, key, value, ttl) error
Delete(ctx, key) error
Increment(ctx, key) (int64, error)
Expire(ctx, key, ttl) error
```

Documented key schema:
- `url:{shortCode}` — TTL 24h (URL data)
- `rate_create:{ip}` — TTL 1h (rate limit counter, not yet implemented)

## .env Notes

- **No inline comments** on value lines — `godotenv` includes them verbatim (breaks int parsing).
- `GOOSE_DBSTRING` must include `port=XXXX` explicitly; it is separate from `DB_PORT`.
- `DB_SSL_MODE=false` maps to `sslmode=disable` in the DSN.

## What Is Incomplete

| Area | Status |
|---|---|
| `UrlRepository` interface + methods | Interface empty, needs full implementation |
| `ClickRepository` | Stub with wrong signatures, not wired in serve.go |
| URL handlers (all 6) | Empty stubs — `CreateShortURL`, `RedirectURL`, `ListMyURLs`, `GetURLDetails`, `DeactivateURL`, `UpdateURL` |
| Redis caching in url_service | Cache injected but never called |
| Rate limiting | `cache.Increment` exists but unused |
| Short code generation | No logic yet (needed by CreateShortURL) |
| `cmd/seed` | Directory exists, main.go empty |
| QR code (`go-qrcode`) | In go.mod, not used anywhere |
| `AdminOnly` middleware | Defined but commented out in router |

## Key File Map

| Path | Purpose |
|---|---|
| `cmd/api/serve.go` | Full dependency wiring |
| `internal/config/config.go` | Singleton config, reads `.env` |
| `internal/db/` | sqlc-generated — do not edit manually |
| `internal/domain/` | Plain domain structs shared across layers |
| `internal/api/middleware/auth_middleware.go` | JWT validation + context keys |
| `internal/response/json.go` | Response helpers + pagination |
| `internal/auth/jwt.go` | Token generation/verification |
| `sql/queries/` | Source SQL for sqlc (edit here, then `make sqlc`) |
| `sql/migrations/` | Goose migration files |
