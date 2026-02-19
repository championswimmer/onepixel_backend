# AGENTS.md

This repository is a Go backend for **onepixel (1px.li)**, an API-first URL shortener with redirect analytics.

## Project purpose and runtime model

- Main entrypoint: `/home/runner/work/onepixel_backend/onepixel_backend/src/main.go`
- Runtime creates three data clients on startup:
  - app DB (`db.GetAppDB`)
  - events DB (`db.GetEventsDB`)
  - GeoIP DB (`db.GetGeoIPDB`)
- HTTP serving uses Fiber with host-based routing:
  - `ADMIN_SITE_HOST` -> admin/API app (`server.CreateAdminApp`)
  - `MAIN_SITE_HOST` -> redirect app (`server.CreateMainApp`)

## Repository map

- `/home/runner/work/onepixel_backend/onepixel_backend/src/main.go`: process startup, host switch, graceful shutdown.
- `/home/runner/work/onepixel_backend/onepixel_backend/src/config`: env loading, constants, config vars.
- `/home/runner/work/onepixel_backend/onepixel_backend/src/server`: Fiber app creation, parser/validator helpers.
- `/home/runner/work/onepixel_backend/onepixel_backend/src/routes`:
  - `api`: `/api/v1/users`, `/api/v1/urls`, `/api/v1/stats`
  - `redirect`: short-code redirects.
- `/home/runner/work/onepixel_backend/onepixel_backend/src/controllers`: business logic for users/urls/events.
- `/home/runner/work/onepixel_backend/onepixel_backend/src/db`: DB providers, singleton initialization, migrations.
- `/home/runner/work/onepixel_backend/onepixel_backend/src/db/models`: GORM models.
- `/home/runner/work/onepixel_backend/onepixel_backend/src/security`: JWT, password hashing, auth middleware.
- `/home/runner/work/onepixel_backend/onepixel_backend/src/dtos`: request/response DTOs.
- `/home/runner/work/onepixel_backend/onepixel_backend/src/utils`: logging, radix64, GeoIP helpers, downloads, PostHog.
- `/home/runner/work/onepixel_backend/onepixel_backend/src/docs`: generated Swagger artifacts.
- `/home/runner/work/onepixel_backend/onepixel_backend/tests`: e2e tests and test DB provider wiring.
- `/home/runner/work/onepixel_backend/onepixel_backend/public_html`: static files served by admin app root.
- `/home/runner/work/onepixel_backend/onepixel_backend/maintenance`: SQL snippets for DB ops.
- `/home/runner/work/onepixel_backend/onepixel_backend/metabase`: Metabase Docker build assets.

## Core behavior details agents should know

### Config and environments

- Config is loaded in `src/config/env.go` with this precedence behavior:
  - If `ENV=test`, load `onepixel.test.env` and chdir to repo root.
  - If production (`ENV=production` or `RAILWAY_ENVIRONMENT=production`), load `onepixel.production.env`.
  - Always load `onepixel.local.env` last as fallback defaults.
- Parsed runtime vars live in `src/config/vars.go`.
- `USE_FILE_DB=true` switches to sqlite/duckdb providers; false uses postgres/clickhouse.

### App/data architecture

- App DB stores users, URL groups, URLs.
- Events DB stores redirect events for analytics.
- DB singletons are initialized via `sync.Once` in `src/db/init.go`.
- Auto-migrations run during DB initialization.
- GeoIP DB file is `./GeoLite2-City.mmdb`; if stale/missing, code downloads it from `https://git.io/GeoLite2-City.mmdb`.

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
- Keep generated Swagger files in sync only when endpoint annotations change.
- Be careful with test/runtime code paths that trigger GeoIP download in network-restricted environments.
- For build and test workflows, see dedicated skills:
  - `/home/runner/work/onepixel_backend/onepixel_backend/.agents/skills/build/SKILL.md`
  - `/home/runner/work/onepixel_backend/onepixel_backend/.agents/skills/test/SKILL.md`
