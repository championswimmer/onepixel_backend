---
name: conventions
description: Apply onepixel repository conventions for env loading, DB backend selection, auth locals, error payloads, and shortcode semantics.
---

# Conventions Skill

Use this skill when touching config, middleware, request/response handling, or URL creation/lookup logic.

## Environment loading order

From `src/config/env.go`:

- `ENV=test` loads `onepixel.test.env` and changes CWD to repo root.
- Production (`ENV=production` or `RAILWAY_ENVIRONMENT=production`) loads `onepixel.production.env`.
- `onepixel.local.env` is loaded as fallback defaults.

## DB backend selection

From `src/config/vars.go` and `src/db/init.go`:

- `USE_FILE_DB=true` -> sqlite + duckdb providers.
- Otherwise -> postgres + clickhouse providers.
- `COMMON_APP_EVENT_DB=true` makes events DB reuse app DB URL.
- Tests override providers via `db.InjectDBProvider` in `tests/providers/test_db_providers.go`.

## Auth locals contract

Use locals keys from `src/config/consts.go`:

- `LOCALS_USER` (`"user"`)
- `LOCALS_ADMIN` (`"admin"`)

## API error response contract

- Use repository helpers (`dtos.CreateErrorResponse`, parser/validator send helpers) instead of ad hoc JSON error payloads.

## Shortcode semantics

- Shortcodes are radix64-encoded numeric IDs (`src/utils/radix64.go`).
- URL creation and lookup logic in `src/controllers/urls.go` assumes radix64 encoding.

## Swagger workflow

- Avoid unnecessary docs regeneration in normal loops.
- Default build path is `make build DOCS=false`.
- Run `make docs` only when route annotations or API docs-facing metadata changed.

