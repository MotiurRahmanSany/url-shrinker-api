# URL Shrinker API

A production-style URL shortener backend built with Go, PostgreSQL, Redis, JWT auth, analytics, and background cleanup.

This project is designed to demonstrate backend engineering fundamentals for real-world systems:
- authentication and session lifecycle
- URL shortening and redirect semantics
- cache-aside Redis strategy
- rate limiting
- click analytics
- graceful shutdown and worker lifecycle

## Highlights

- REST API with clean layered architecture: handler -> service -> repository -> sqlc
- JWT access tokens + refresh token rotation
- Short URL creation with optional custom code
- Public redirect endpoint with 302 behavior
- Expiration, deactivation, and max-click protection
- Redis cache for redirect hot path
- Redis-based rate limiting for URL creation
- Click tracking and URL analytics endpoint
- Hourly background worker to delete expired URLs

## Architecture

```text
HTTP Handlers -> Services -> Repositories -> sqlc Queries -> PostgreSQL
										 |
										 v
									 Redis

Background Worker -> Repository -> sqlc -> PostgreSQL
```

Main wiring is in cmd/api/serve.go.

## Tech Stack

- Go 1.25
- net/http (standard library router/mux)
- PostgreSQL 16
- pgx/v5
- sqlc
- goose
- Redis 7 (go-redis/v9)
- JWT (github.com/golang-jwt/jwt/v5)
- bcrypt (golang.org/x/crypto)
- Docker Compose

## Project Structure

```text
cmd/
	api/                # API bootstrap and dependency wiring
	seed/               # admin seeding command
internal/
	api/
		handlers/
		middleware/
		router/
	auth/
	cache/
	config/
	database/
	db/                 # sqlc-generated code
	domain/
	repository/
	response/
	service/
	worker/
sql/
	migrations/
	queries/
```

## Prerequisites

- Go 1.25+
- Docker + Docker Compose
- goose CLI
- sqlc CLI (only needed when query files change)
- air (optional, used by make run)

## Environment Variables

Important: .env is required. The app intentionally fails fast if .env is missing.

Use these values in your .env file:

| Variable | Required | Example | Purpose |
|---|---|---|---|
| APP_ENV | no | development | logging mode (text/json) |
| VERSION | yes (recommended) | 1.0.0 | service version |
| HTTP_PORT | yes (recommended) | 8080 | API port |
| SERVICE_NAME | yes (recommended) | URL Shrinker API | service name |
| JWT_SECRET_KEY | yes | super_secret_key | JWT signing key |
| DB_HOST | yes | localhost | PostgreSQL host |
| DB_PORT | yes | 5432 | PostgreSQL port |
| DB_USER | yes | postgres | PostgreSQL user |
| DB_PASSWORD | yes | postgres | PostgreSQL password |
| DB_NAME | yes | url_shrinker | PostgreSQL database |
| DB_SSL_MODE | yes | false | SSL toggle |
| REDIS_HOST | yes | localhost | Redis host |
| REDIS_PORT | yes | 6379 | Redis port |
| REDIS_PASSWORD | no |  | Redis password |
| REDIS_DB | yes | 0 | Redis DB index |
| GOOSE_DBSTRING | yes (for migrations) | host=... port=... user=... password=... dbname=... sslmode=disable | goose DSN |
| SEED_ADMIN_EMAIL | no (only for seed-admin) | admin@example.com | admin seed email |
| SEED_ADMIN_PASSWORD | no (only for seed-admin) | admin123456 | admin seed password |

## Local Run

1. Start infrastructure (Postgres + Redis)

```bash
docker compose up -d
```

2. Run migrations

```bash
make migrate-up
```

3. Start API

Option A (recommended while developing):

```bash
make run
```

Option B (without air):

```bash
go run ./cmd/api
```

4. Stop services

```bash
make stop
```

## Useful Commands

```bash
make sqlc          # regenerate sqlc output from sql/queries
make migrate-up    # apply migrations
make migrate-down  # rollback migration
make flush-redis   # clear redis data
make seed-admin    # create admin user if not present
```

## API Overview

Base URL:

```text
http://localhost:8080
```

### Public Endpoints

- GET /health
- POST /auth/register
- POST /auth/login
- POST /auth/refresh
- GET /{code}

### Protected Endpoints (Authorization: Bearer <token>)

- GET /auth/me
- POST /auth/logout
- POST /urls
- GET /urls
- GET /urls/{code}
- PATCH /urls/{code}
- DELETE /urls/{code}
- GET /urls/{code}/stats

## Response Envelope

All responses use a common shape:

```json
{
	"success": true,
	"message": "...",
	"data": {}
}
```

Error example:

```json
{
	"success": false,
	"message": "...",
	"error": "..."
}
```

## Quick API Flow (curl)

Register:

```bash
curl -X POST http://localhost:8080/auth/register \
	-H "Content-Type: application/json" \
	-d '{"email":"user@example.com","password":"password123"}'
```

Login:

```bash
curl -X POST http://localhost:8080/auth/login \
	-H "Content-Type: application/json" \
	-d '{"email":"user@example.com","password":"password123"}'
```

Create URL (replace ACCESS_TOKEN):

```bash
curl -X POST http://localhost:8080/urls \
	-H "Authorization: Bearer ACCESS_TOKEN" \
	-H "Content-Type: application/json" \
	-d '{"original_url":"https://example.com","custom_short_code":"mycode123"}'
```

Redirect:

```bash
curl -i http://localhost:8080/mycode123
```

Get stats:

```bash
curl -X GET http://localhost:8080/urls/mycode123/stats \
	-H "Authorization: Bearer ACCESS_TOKEN"
```

## Behavior Semantics

- 302 Found for redirect endpoint
- 410 Gone when URL is inactive, expired, or max-click limit is reached
- 429 Too Many Requests when create-URL rate limit is exceeded
- Redis failures in rate limiting are fail-open (request is allowed)
- Redis failures in redirect cache degrade gracefully to database lookup

## Data Model

Core tables:
- users
- refresh_tokens
- urls
- clicks

See migration files in sql/migrations for full schema.

## Background Worker

- Worker runs hourly
- Deletes expired URLs from database
- Starts with API server and stops on shutdown context cancellation

## Security Notes

- Access tokens are JWTs
- Refresh tokens are random opaque values stored in database
- Passwords are hashed with bcrypt
- Ownership checks are enforced for user URL operations

## Current Limitations

- No automated test suite yet (unit/integration tests pending)
- Pagination total count for GET /urls is currently based on returned slice length, not a dedicated COUNT query
- Admin-only routes exist as middleware but are not yet exposed in router
- Some handler comments are still verbose and can be cleaned up

## Suggested Next Steps

1. Add integration tests for auth + create + redirect + stats.
2. Add COUNT query for accurate pagination metadata.
3. Introduce request ID middleware and structured request logging fields.
4. Add optional admin endpoints and role-based router groups.
5. Add CI workflow for lint + test + migration checks.

## Author

Motiur Rahman Sany

