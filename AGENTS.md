# AGENTS.md

This repository is a Go backend for **onepixel (1px.li)**, an API-first URL shortener with redirect analytics.

## Project purpose and runtime model

- Main entrypoint: `src/main.go`
- Runtime creates three data clients on startup:
  - app DB (`db.GetAppDB`)
  - events DB (`db.GetEventsDB`)
  - GeoIP DB (`db.GetGeoIPDB`)
- HTTP serving uses Fiber with host-based routing:
  - `ADMIN_SITE_HOST` -> admin/API app (`server.CreateAdminApp`)
  - `MAIN_SITE_HOST` -> redirect app (`server.CreateMainApp`)

## Build and test commands

- Fast build (skip Swagger regeneration): `make build DOCS=false`
- Full build with Swagger: `make build` (requires `swag` installed via `go install github.com/swaggo/swag/cmd/swag@latest`)
- Full multi-OS release artifacts: `make build_all DOCS=false`
- Regenerate Swagger docs (only when API annotations/routes change): `make docs`
- Full test suite: `make test`
- Unit tests only: `make test_unit`
- E2E tests only: `make test_e2e`
- Single e2e test: `ENV=test go test ./tests/e2e -run TestRegisterUser -count=1 -v`
- Single unit test: `ENV=test go test ./src/controllers/... -run TestCreateRandomShortUrl -count=1 -v`
- `make test_clean` removes `app.db` and `events.db` before each test run automatically.

## Repository map

- `src/main.go`: process startup, host switch, graceful shutdown.
- `src/config`: env loading, constants, config vars.
- `src/server`: Fiber app creation, parser/validator helpers.
- `src/routes`:
  - `api`: `/api/v1/users`, `/api/v1/urls`, `/api/v1/stats`
  - `redirect`: short-code redirects.
- `src/controllers`: business logic for users/urls/events.
- `src/db`: DB providers, singleton initialization, migrations.
- `src/db/models`: GORM models.
- `src/security`: JWT, password hashing, auth middleware.
- `src/dtos`: request/response DTOs.
- `src/utils`: logging, radix64, GeoIP helpers, downloads, PostHog.
- `src/docs`: generated Swagger artifacts.
- `tests`: e2e tests and test DB provider wiring.
- `tests/providers/test_db_providers.go`: overrides DB providers for test isolation via `db.InjectDBProvider`.
- `public_html`: static files served by admin app root.
- `maintenance`: SQL snippets for DB ops.
- `metabase`: Metabase Docker build assets.

## Core behavior details agents should know

### Config and environments

- Config is loaded in `src/config/env.go` with this precedence behavior:
  - If `ENV=test`, load `onepixel.test.env` and chdir to repo root.
  - If production (`ENV=production` or `RAILWAY_ENVIRONMENT=production`), load `onepixel.production.env`.
  - Always load `onepixel.local.env` last as fallback defaults.
- Parsed runtime vars live in `src/config/vars.go`.
- `USE_FILE_DB=true` switches to sqlite/duckdb providers; false uses postgres/clickhouse.
- `COMMON_APP_EVENT_DB=true` makes the events DB reuse the app DB URL (useful for single-file local dev).

### App/data architecture

- App DB stores users, URL groups, URLs.
- Events DB stores redirect events for analytics.
- DB singletons are initialized via `sync.Once` in `src/db/init.go`.
- Auto-migrations run during DB initialization.
- GeoIP DB file is `./GeoLite2-City.mmdb`; if stale/missing, code downloads it from `https://git.io/GeoLite2-City.mmdb`.
- Shortcodes are radix64-encoded numeric IDs (`src/utils/radix64.go`); random URLs use a 6-char radix64 space (64^6 ≈ 68B entries). Lookup and creation logic in `src/controllers/urls.go` assumes this encoding.
- DB provider injection (`db.InjectDBProvider(name, fn)`) allows tests to substitute sqlite/duckdb without changing production code; see `tests/providers/test_db_providers.go`.

### HTTP/API architecture

- Admin app (`src/server/server.go`):
  - `/api/v1/users` (registration/login/update)
  - `/api/v1/urls` (create/query URLs)
  - `/api/v1/stats` (event stats)
  - `/docs/*` (Swagger UI)
  - static web pages from `public_html`
- Main app (`src/server/server.go` + `src/routes/redirect/redirect.go`):
  - redirects by shortcode
  - logs redirect events asynchronously
- Route handler layering: route handler → `server/parsers` (body parse) → `server/validators` (field validation) → controller method → DTO response. Use `dtos.CreateErrorResponse` for all error payloads.

### Auth and security

- JWT auth uses `Authorization: Bearer <token>`.
- Admin API key uses `X-API-Key` and `config.AdminApiKey`.
- Middleware stores values in Fiber locals with keys from `src/config/consts.go`:
  - `user`
  - `admin`
- Passwords are hashed using bcrypt in `src/security/password_hash.go`.

### Analytics and external integrations

- Redirect logging enriches events with GeoIP data.
- PostHog events are emitted via `src/utils/posthog/posthog_client.go`.

## Working guidance for future agents

- Prefer minimal, surgical changes; preserve existing public behavior.
- Keep generated Swagger files in sync only when endpoint annotations change; skip with `DOCS=false` otherwise.
- Be careful with test/runtime code paths that trigger GeoIP download in network-restricted environments.
- Use `make build DOCS=false` as the default fast build path.

## Build and lint notes

- No dedicated lint target is configured in `Makefile` or CI workflow.
