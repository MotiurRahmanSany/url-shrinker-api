# URL Shrinker API - Full End-to-End Roadmap

## Current Status (aligned with Q&A blueprint)

| Phase | Status | Notes |
|-------|--------|-------|
| Phase 1: Foundation | Done | Docker, config, DB pool, Redis, logger middleware |
| Phase 2: Auth | Done | Register, login, refresh, logout, me working |
| Phase 3: URL Core + Redirect | In Progress | Migrations + sqlc done; handlers/services/repositories still incomplete |
| Phase 4: Redis Cache | Not started | Cache wired but not used in redirect flow |
| Phase 5: Rate Limiting | Not started | Redis increment exists but middleware not implemented |
| Phase 6: Click Analytics | Not started | Click repository/service/handler incomplete |
| Phase 7: Background Job | Not started | Cleanup worker not implemented |
| Phase 8: Optional Enhancements | Not started | Preview/QR/Search/Admin pending |
| Phase 9: Polish + README | Not started | Docs/tests/final hardening pending |

### Immediate Blocking Issues (fix first)
1. `internal/repository/click_repository.go`: `CreateClick` uses undefined `Click` type (must use `domain.Click`).
2. `internal/repository/url_repository.go`: interface is empty; add method signatures and implementations.
3. `cmd/api/serve.go`: `ClickRepository` not wired into service graph.

---

## 1. Project Goal
Build a production-style backend URL shortener for junior backend portfolio depth: auth, URL CRUD, public redirect, analytics, Redis caching, rate limiting, and cleanup worker.

## 2. Architecture

```text
HTTP Handlers -> Services -> Repositories -> sqlc Queries -> PostgreSQL
                    |
                    v
              Cache Interface -> Redis
```

- Handlers: HTTP only (decode request, call service, map errors, return response).
- Services: business rules, sentinel errors, ownership checks, cache strategy.
- Repositories: DB-only, map pgtype values to domain types.
- Cache: cache-aside strategy for redirect hot path.
- Dependency wiring: `cmd/api/serve.go`.

## 3. Tech Stack
- Go (`net/http`)
- PostgreSQL + `pgx/v5`
- `sqlc`
- `goose`
- Redis + `go-redis/v9`
- JWT (`golang-jwt/jwt/v5`)
- bcrypt (`golang.org/x/crypto`)
- `godotenv`
- Docker Compose
- Optional: `go-qrcode`

## 4. Database Schema

### users
- `id UUID PK`
- `email TEXT UNIQUE NOT NULL`
- `password_hash TEXT NOT NULL`
- `role TEXT CHECK('admin','student')`
- `is_active BOOLEAN DEFAULT true`
- `created_at TIMESTAMP`
- `updated_at TIMESTAMP`

### refresh_tokens
- `id BIGSERIAL PK`
- `user_id UUID FK -> users.id ON DELETE CASCADE`
- `token TEXT UNIQUE`
- `expires_at TIMESTAMP`
- `revoked BOOLEAN DEFAULT false`
- `created_at TIMESTAMP`

### urls
- `id BIGSERIAL PK`
- `short_code VARCHAR(10) UNIQUE NOT NULL`
- `original_url TEXT NOT NULL`
- `user_id UUID FK -> users.id ON DELETE SET NULL`
- `is_active BOOLEAN NOT NULL DEFAULT true`
- `expires_at TIMESTAMPTZ NULL`
- `max_clicks INTEGER NULL`
- `created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()`
- `updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()`

### clicks
- `id BIGSERIAL PK`
- `url_id BIGINT NOT NULL FK -> urls.id ON DELETE CASCADE`
- `clicked_at TIMESTAMPTZ NOT NULL DEFAULT NOW()`
- `ip_address VARCHAR(45)`
- `user_agent TEXT`
- `referer TEXT`

### Required Indexes
- `urls(short_code)`
- `urls(user_id)`
- `clicks(url_id)`

## 5. API Endpoints

### Public (no auth)
| Method | Path | Purpose | Status |
|--------|------|---------|--------|
| GET | `/health` | Health check | Done |
| POST | `/auth/register` | Register | Done |
| POST | `/auth/login` | Login | Done |
| POST | `/auth/refresh` | Refresh token | Done |
| GET | `/{code}` | 302 redirect to original URL | Stub |

### Protected (JWT)
| Method | Path | Purpose | Status |
|--------|------|---------|--------|
| GET | `/auth/me` | Current user | Done |
| POST | `/auth/logout` | Logout | Done |
| POST | `/urls` | Create short URL | Stub |
| GET | `/urls` | List my URLs | Stub |
| GET | `/urls/{code}` | URL details | Stub |
| PATCH | `/urls/{code}` | Update URL | Stub |
| DELETE | `/urls/{code}` | Soft delete URL | Stub |
| GET | `/urls/{code}/stats` | Analytics | Not created |

### Optional
| Method | Path | Purpose |
|--------|------|---------|
| GET | `/{code}/preview` | Show destination without redirect |
| GET | `/urls/{code}/qr` | Generate QR code image |
| GET | `/admin/urls` | Admin list all URLs |
| DELETE | `/admin/urls/{code}` | Admin force deactivate |

Router correction required:
- Redirect should be `GET /{code}` (public).
- Details should be `GET /urls/{code}` (protected).
- Use `PATCH` for update.
- Add `GET /urls/{code}/stats`.

## 6. Core Features (must-have)
1. Auth lifecycle (register/login/refresh/logout/me).
2. Create short URL (auth required).
3. Optional custom short code.
4. Redirect with `302 Found`.
5. Expired/inactive links return `410 Gone`.
6. Optional max-click limit enforcement.
7. Redis cache for redirect path.
8. Rate limiting on create endpoint (e.g., 10/hour/IP).
9. List/details/update/deactivate URL.
10. Click analytics (total, today, per-day timeline).
11. Cleanup worker for expired URL hard delete.

## 7. Scope Discipline
Avoid these before core completion:
- Custom domains.
- Frontend dashboard.
- Rich analytics UI.
- Link preview metadata scraping.

Keep optional features only after Phase 7 is stable.

## 8. Business Rules
- Accept only `http://` or `https://` as original URL.
- `short_code` must be unique (DB constraint is final guard).
- Ownership checks required for details/update/delete/stats.
- Soft delete sets `is_active=false`.
- Redirect returns `410` when inactive, expired, or max clicks reached.
- Redis failure must degrade gracefully to DB.
- Pagination with `page` and `limit`; include `total` and `total_pages`.

## 9. Redis Plan

### Redirect cache (cache-aside)
- Key: `url:{short_code}`
- TTL: 24h (or clamp to `expires_at`)
- Read flow: cache -> DB on miss -> write cache
- Invalidate on update/deactivate: `DEL url:{short_code}`

### Rate limiting
- Key: `rate_create:{ip}`
- Pattern: `INCR` then `EXPIRE` on first hit
- Window: 1h
- Limit: 10 requests per hour per IP
- Exceeded: `429 Too Many Requests`

## 10. Error to HTTP Mapping
- Validation error -> `400 Bad Request`
- Invalid creds/JWT -> `401 Unauthorized`
- Forbidden owner action -> `403 Forbidden`
- Not found -> `404 Not Found`
- Conflict (email/code taken) -> `409 Conflict`
- Expired/inactive/max clicks -> `410 Gone`
- Rate limited -> `429 Too Many Requests`
- Internal/database -> `500 Internal Server Error`

## 11. Implementation Phases (from Q&A blueprint)

### Phase 1 - Foundation (2-3h) - Done
- Environment, Docker, config, DB pool, Redis client, base middleware/response.

### Phase 2 - Auth (1-2h) - Done
- User/token queries, repositories, service, handlers, middleware, routes.

### Phase 3 - URL Core + Redirect (4-5h) - In Progress
- Fix repository blockers (`UrlRepository`, `ClickRepository`, DI wiring).
- Base62 generator with collision retry.
- URL validation.
- Implement handlers: create/list/details/update/deactivate.
- Add redirect `GET /{code}` with `302`.
- Enforce `410` for inactive/expired/max-clicks.

### Phase 4 - Redis Cache (2-3h)
- Redirect path reads from cache first.
- Miss falls back to DB and repopulates cache.
- Invalidate cache on update/deactivate.
- Do not fail redirect if Redis is down.

### Phase 5 - Rate Limiting (2-3h)
- Middleware for `POST /urls` using Redis `INCR/EXPIRE`.
- Return `429` using common response envelope.

### Phase 6 - Click Analytics (3-4h)
- Complete click repository methods.
- Record click metadata on successful redirect.
- Add stats service + handler: total, today, daily timeline.

### Phase 7 - Background Job (1-2h)
- Hourly ticker worker.
- Delete expired URLs (`DELETE WHERE expires_at < NOW()`).
- Graceful shutdown with context cancellation.

### Phase 8 - Optional Enhancements (2-3h)
- Preview endpoint.
- QR endpoint.
- Search/filter list endpoint.
- Optional admin moderation endpoints.

### Phase 9 - Polish + README (2h)
- README with setup + endpoint examples.
- Seed/admin improvements.
- Error consistency pass and ownership/security checks.
- Basic integration tests for critical flows.

## 12. Time Tracking Summary

| Phase | Estimate |
|------|----------|
| 1 | 2-3h |
| 2 | 1-2h |
| 3 | 4-5h |
| 4 | 2-3h |
| 5 | 2-3h |
| 6 | 3-4h |
| 7 | 1-2h |
| 8 | 2-3h |
| 9 | 2h |
| Total | 20-27h |

Recommended pace: 6-8 days with verification after each phase.

## 13. Expected Learning Outcomes
- Clean architecture in Go without framework.
- SQL-first workflow with `sqlc` and migrations.
- JWT + refresh token rotation.
- `pgtype` to domain mapping with `pgx/v5`.
- Redis cache-aside and invalidation strategy.
- Redis rate limiting (`INCR` + `EXPIRE`).
- Correct HTTP semantics (`302`, `410`, `429`).
- Analytics query design with aggregation.
- Background workers and graceful shutdown.
- Production-oriented API and security decisions.

## 14. Interview Prep Prompts
- Why Base62 random over UUID?
- Why `302` instead of `301`?
- How cache invalidation works on update/deactivate?
- What happens when Redis is down?
- How abuse prevention works?
- How this design scales under high redirect traffic?
