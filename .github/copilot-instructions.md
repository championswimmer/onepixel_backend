# Copilot Instructions for `onepixel_backend`

## Build, test, and lint commands

### Build

- Preferred build (skip Swagger regeneration):
  - `make build DOCS=false`
- Full multi-OS release artifacts:
  - `make build_all DOCS=false`
- Regenerate Swagger docs (only when API annotations/routes change):
  - `make docs`

### Tests

- Full test suite:
  - `make test`
- Unit tests only:
  - `make test_unit`
- E2E tests only:
  - `make test_e2e`
- Single test (example):
  - `ENV=test go test ./tests/e2e -run TestRegisterUser -count=1 -v`
- Single package/unit target (example):
  - `ENV=test go test ./src/controllers/... -run TestCreateRandomShortUrl -count=1 -v`

### Lint

- No dedicated lint target is configured in `Makefile` or CI workflow.

## High-level architecture

- The process starts in `src/main.go` and initializes three singleton data clients:
  - application DB (`db.GetAppDB`)
  - events DB (`db.GetEventsDB`)
  - GeoIP DB (`db.GetGeoIPDB`)
- HTTP is served by a top-level Fiber app that dispatches by `Host` header:
  - `ADMIN_SITE_HOST` -> admin/API app (`server.CreateAdminApp`)
  - `MAIN_SITE_HOST` -> redirect app (`server.CreateMainApp`)
- Data is split by responsibility:
  - app DB stores users, URL groups, and URLs
  - events DB stores redirect analytics events
- Redirect flow (`src/routes/redirect/redirect.go`) resolves shortcode -> redirects -> logs analytics asynchronously via `EventsController.LogRedirectAsync` (GeoIP enrichment + PostHog emit).
- API routing uses a layered pattern:
  - route handler -> parser (`server/parsers`) -> validator (`server/validators`) -> controller -> DTO response

## Key conventions for this repository

- Env loading order in `src/config/env.go` is important:
  - `ENV=test` loads `onepixel.test.env` and changes CWD to repo root
  - production loads `onepixel.production.env`
  - `onepixel.local.env` is always loaded as fallback defaults
- DB backend selection is config-driven:
  - `USE_FILE_DB=true` uses sqlite/duckdb providers
  - otherwise postgres/clickhouse providers are used
  - tests inject their own providers in `tests/providers/test_db_providers.go`
- Auth context is passed through Fiber locals using constants from `src/config/consts.go`:
  - `LOCALS_USER` (`"user"`)
  - `LOCALS_ADMIN` (`"admin"`)
- API errors should use repository DTO/error helpers (`dtos.CreateErrorResponse`, parser/validator error helpers) instead of ad hoc payloads.
- Shortcodes are radix64-encoded numeric IDs (`src/utils/radix64.go`, `src/controllers/urls.go`); URL lookup/creation logic assumes this encoding.
- Avoid unnecessary Swagger regeneration in normal coding loops; `make build DOCS=false` is the default fast path, and docs are generated from route annotations in `src/server/server.go` and route files.
