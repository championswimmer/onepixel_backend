# Test Skill

Use this skill for unit/e2e test execution, targeted package tests, and test troubleshooting.

## Primary test targets

Defined in `/home/runner/work/onepixel_backend/onepixel_backend/Makefile`:

- `make test_unit`
  - Cleans `app.db`, `events.db`, `events.db.wal`
  - Runs `go test` over `./src/...` with race + coverage output `coverage.unit.out`
- `make test_e2e`
  - Cleans DB files
  - Runs `go test` over `./tests/...` with race + coverage output `coverage.e2e.out`
- `make test`
  - Runs `test_unit` then `test_e2e`

## Test environment mechanics

- Tests are expected to run with `ENV=test` (set by Makefile test targets).
- `src/config/env.go` under `ENV=test`:
  - changes directory to project root
  - loads `onepixel.test.env`
- `tests/providers/test_db_providers.go` injects sqlite + duckdb providers for tests.

## Recommended execution flow for agents

1. Run targeted tests for touched package(s) first.
2. Then run broader target (`make test_unit` or `make test`) as needed.
3. Keep changes minimal and avoid altering unrelated flaky tests.

## Known failure mode in restricted environments

- Some tests initialize `db.GetGeoIPDB()` which may attempt to download:
  - `https://git.io/GeoLite2-City.mmdb`
- In network-restricted environments this can fail with DNS/network errors and panic.
- If this occurs, treat it as environment-related unless your change touched GeoIP behavior.

## Useful targeted commands

- Package-level unit tests:
  - `go test ./src/security/...`
  - `go test ./src/controllers/...`
- Single test run example:
  - `go test ./tests/e2e -run TestRegisterUser`

## Test artifacts and cleanup

- Coverage outputs:
  - `coverage.unit.out`
  - `coverage.e2e.out`
- DB files used/cleaned:
  - `app.db`, `events.db`, `events.db.wal`
