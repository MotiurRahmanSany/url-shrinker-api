# URL Shrinker API

A production-grade URL shortener backend demonstrating robust backend engineering fundamentals. Built with **Go**, **PostgreSQL**, **Redis**, **JWT auth**, **analytics**, and **background workers**.

## Why This Project?

URL shortening is deceptively simple on the surface but reveals real engineering depth:
- **High read throughput** — redirects happen constantly.
- **Cache strategy** — trade-offs between Redis and DB.
- **Rate limiting** — prevent abuse.
- **Analytics** — track usage patterns.
- **Proper HTTP semantics** — 302 vs 301, 410 Gone states.
- **Background jobs** — housekeeping at scale.

This project is built with intentional architecture decisions and production awareness.

## Features

### Core
- ✅ **Authentication**: Secure JWT access tokens with refresh token rotation.
- ✅ **Custom Short Codes**: Create short URLs with random or user-defined custom codes.
- ✅ **High-Performance Redirects**: Public 302 redirect endpoint backed by Redis cache (`cache-aside` strategy).
- ✅ **Safety Controls**: Expiration limits, deactivation, and max-click protection (`410 Gone`).
- ✅ **Analytics Engine**: Granular click tracking and insights endpoints.
- ✅ **Rate Limiting**: Redis-based sliding window rate limiter (e.g., 10 requests/hour/IP).
- ✅ **Garbage Collection**: Hourly background worker routine for expired URL cleanup.
- ✅ **Graceful Shutdown**: Context cancellation blocks ensuring safe terminations.

### Tech Stack

| Component | Technology |
|-----------|-----------|
| **Language** | Go 1.25+ |
| **HTTP Router** | standard `net/http` |
| **Database** | PostgreSQL 16 |
| **Migrations** | `goose` |
| **SQL Engine** | `sqlc` + `pgx/v5` driver |
| **Cache** | Redis 7 (`go-redis/v9`)|
| **Cryptography** | `bcrypt` + `golang-jwt/v5` |
| **Orchestration**| Docker Compose |

---

## Architecture Breakdown

```text
HTTP Request
     ↓
Logger Middleware → CORS Middleware
     ↓
Auth / Rate Limit Middleware (if protected)
     ↓
Handler (decodes request)
     ↓
Service (business logic, sentinel errors, domain validations)
     ↓
Repository (DB mapping, pgtype conversions)
     ↓
sqlc Queries ← → PostgreSQL
     ↓
Cache Layer ← → Redis
     ↓
HTTP Response
```

**Wiring & Dependency Injection**: Centralized clearly in `cmd/api/serve.go`.

---

## Complete Project Structure

```text
cmd/
├── api/             # Entry point and dependency injection wiring
└── seed/            # Admin user seeding scripts

internal/
├── api/             
│   ├── handlers/    # Payload mapping, context assignments (Auth, URL, Clicks)
│   ├── middleware/  # CORS, JWT protections, Rate Limiting logic
│   └── router/      # HTTP stdlib mux mappings
├── auth/            # JWT signing & validation utilities
├── cache/           # Redis interfaces and abstractions
├── config/          # .env processing parameters
├── database/        # pgx database pooling setup
├── db/              # Auto-generated code from sqlc
├── domain/          # Plain domain objects / Entities mapping logic
├── repository/      # Strict repository interfaces bridging DB/Cache
├── response/        # Consistent HTTP response enveloping
├── service/         # Domain-level business rules avoiding HTTP concerns
└── worker/          # Background ticker routines (Cleanup Worker)

client/              # Next.js React frontend 
sql/
├── migrations/      # Goose version-controlled schema increments
└── queries/         # Original SQL queries mapped to sqlc Go interfaces
```

---

## Database & Data Flow

- The structure isolates external packages. `Handlers` process `json` bodies and push mapping interfaces to `Services` handling structural checks against `domain` payloads.
- `Repositories` shield services from SQL-types (`pgtype`) handling mappings exclusively into pure Go structs.

### Data Models
- `users`: Core authentication identity mapping.
- `refresh_tokens`: Expiry-managed relational rows enabling session retention.
- `urls`: Maps `custom_short_code` target resolutions containing constraints (`expires_at`, `max_clicks`).
- `clicks`: Aggregated append-only ledger processing request footprints (`User-Agent`, `Referer`, `IP`).

---

## API Endpoints

### Public Endpoints

| Method | Path | Semantics |
|--------|------|-----------|
| GET | `/health` | Sanity checks and load-balancer probes |
| POST | `/auth/register` | Open account initiation |
| POST | `/auth/login` | Bearer payload exchanges |
| POST | `/auth/refresh` | Revise TTL windows using valid payload IDs |
| GET | `/{code}` | **302 Found** Target resolver (Redis Hotpath) |

*Notes on Redirect:* If a target exceeds `max_clicks`, is deactivated, or expires manually, the route cleanly resolves HTTP `410 Gone`.

### Protected Endpoints (`Authorization: Bearer <token>`)

| Method | Path | Semantics |
|--------|------|-----------|
| POST | `/urls` | Generate new mapping rules. Rate limited dynamically by `REDIS` bounds returning `429 Too Many Requests`. |
| GET | `/urls` | Paginated lists indexing matching `user_id` blocks |
| GET | `/urls/{code}`| View full state properties |
| PATCH | `/urls/{code}`| Update temporal rules (expiration, max metrics) |
| DELETE | `/urls/{code}`| Toggle `is_active` constraints natively terminating routing paths |
| GET | `/urls/{code}/stats`| Combine DB constraints against the `/clicks` aggregate block arrays |

---

## Quick Start Development

### Prerequisites
- Go 1.25+
- Docker + Docker Compose
- `goose` CLI
- `sqlc` CLI (if you intend on running custom modification generators)

### Installation & Launch

1. Clone and setup configurations:
```bash
cp .env.example .env
```
Ensure your `.env` is configured correctly (especially the `GOOSE_DBSTRING` for migrations):

```env
APP_ENV=development
HTTP_PORT=8080
SERVICE_NAME=URL Shrinker
JWT_SECRET_KEY=your_super_secret_jwt_key

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=url_shrinker_db
DB_SSL_MODE=false

REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

GOOSE_DRIVER=postgres
GOOSE_DBSTRING=host=localhost port=5432 user=postgres password=postgres dbname=url_shrinker_db sslmode=disable
GOOSE_MIGRATION_DIR=./sql/migrations
```

2. Execute Docker infrastructures:
```bash
docker compose up -d
```

3. Populate schemas using goose:
```bash
make migrate-up
```

4. Run the live reloading compiler server (`air`):
```bash
make run
```
Your backend will start actively processing at `http://localhost:8080/health`.

### Testing

Generate a shortlink (Requires JWT token from Login step):
```bash
curl -X POST http://localhost:8080/urls \
	-H "Authorization: Bearer ACCESS_TOKEN" \
	-H "Content-Type: application/json" \
	-d '{"original_url":"https://github.com", "custom_short_code":"gh"}'
```

Verify Redirect mapping (Fast Redis execution):
```bash
curl -i http://localhost:8080/gh
```

### Automation Scripts (`Makefile`)

| Command | Action |
|---------|--------|
| `make sqlc` | Regenerate mappings reflecting changes in `sql/queries`. |
| `make migrate-up` | Run goose migrations reflecting structural definitions. |
| `make migrate-down` | Rollback single block structural layers. |
| `make flush-redis` | Clear local active cache bindings. |
| `make seed-admin` | Trigger a pre-configured domain creation entity. |

