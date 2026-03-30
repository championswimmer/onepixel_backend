---
name: database
description: Understand onepixel database topology by environment, data ownership across app/events stores, and initialization/migration behavior.
---

# Database Skill

Use this skill when changing config, models, controllers, or test setup that touches persistence.

## Which databases are used in each environment

- Local/dev defaults (`onepixel.local.env`):
  - App DB: `sqlite` (`DB_DIALECT=sqlite`, `DATABASE_URL=app.db`)
  - Events DB: `duckdb` (`EVENTDB_DIALECT=duckdb`, `EVENTDB_URL=events.db`)
  - `USE_FILE_DB=true` enables file-backed providers.
- Test (`ENV=test` + `onepixel.test.env`):
  - App DB: `sqlite` (`app.db`)
  - Events DB: `duckdb` (`events.db`)
  - Tests import `onepixel_backend/tests/providers` to inject test DB providers.
- Production template (`onepixel.production.env`):
  - `USE_FILE_DB=false` enables network DB providers.
  - App DB: `postgres` (`DB_DIALECT=postgres`).
  - Events DB: also `postgres` by default because `COMMON_APP_EVENT_DB=true` and `EVENTDB_DIALECT=postgres`.
  - Note: code supports `clickhouse` for events when `EVENTDB_DIALECT=clickhouse` and `COMMON_APP_EVENT_DB=false`.

## What data is stored in which database

- App DB (OLTP/domain entities):
  - `users` (`models.User`)
  - `url_groups` (`models.UrlGroup`)
  - `urls` (`models.Url`)
- Events DB (analytics/time-series style):
  - `events_redirect` (`models.EventRedirect`)
  - Stores redirect telemetry: short URL identifiers, creator/group IDs, IP/user-agent, referer, and GeoIP fields.
- GeoIP DB:
  - Separate MMDB file (`./GeoLite2-City.mmdb`), opened by `db.GetGeoIPDB()` for enrichment; not managed by GORM migrations.

## How DBs are initialized

- Initialization is singleton-based in `src/db/init.go`:
  - `GetAppDB`, `GetEventsDB`, and `GetGeoIPDB` each guarded by `sync.Once`.
- Runtime startup in `src/main.go` eagerly initializes all three DB clients.
- Providers are registered in `db.init()` based on `config.UseFileDB`:
  - file mode: `sqlite`, `duckdb`
  - network mode: `postgres`, `clickhouse`
- Provider implementations live in `src/db/providers.go`:
  - `postgres`/`clickhouse` open with retry (`attemptToOpen`, 10 attempts with 1s delay).
  - `sqlite`/`duckdb` open directly.

## How schema changes (migrations) work

- Migrations are automatic via GORM `AutoMigrate` in `src/db/init.go`.
- App DB migration flow (`GetAppDB`):
  - Migrates `User`, `UrlGroup`, `Url`.
  - Failures panic (`lo.Must0`), so startup/first access fails fast.
- Events DB migration flow (`GetEventsDB`):
  - Attempts to migrate `EventRedirect`.
  - Failure is logged but not fatal (`lo.TryWithErrorValue` + `applogger.Error`).
- There are no SQL migration files or versioned migration runner; model changes are applied opportunistically at DB initialization time.

## Test-specific DB wiring details

- `tests/providers/test_db_providers.go` overrides `sqlite` and `duckdb` providers via `db.InjectDBProvider`.
- Tests activate overrides through side-effect imports (`_ "onepixel_backend/tests/providers"`).
- `make test_clean` removes `app.db` and `events.db` before each test suite target.
- `ENV=test` also changes cwd to repo root in `src/config/env.go`, ensuring relative DB paths resolve consistently.

