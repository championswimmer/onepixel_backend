---
name: test
description: Run onepixel unit and e2e tests with correct env setup, coverage artifacts, and troubleshooting guidance for GeoIP/network-related failures.
---

# Test Skill

Use this skill for unit tests, e2e tests, targeted test runs, and test triage.

## Quick commands

- Full unit suite:
  - `make test_unit`
- Full e2e suite:
  - `make test_e2e`
- Full test pipeline:
  - `make test`
- Targeted package tests (fast feedback):
  - `ENV=test go test -count 1 ./src/security/...`
  - `ENV=test go test -count 1 ./src/controllers/...`
- Single e2e test:
  - `ENV=test go test -count 1 ./tests/e2e -run TestRegisterUser`

## Current Makefile test behavior

Defined in `Makefile` at repo root:

- `test_clean`
  - removes `app.db`, `events.db`, `events.db.wal`
- `test_unit` (depends on `test_clean`)
  - removes `coverage.unit.out`
  - runs:
    - `ENV=test go test -count 1 -timeout 10s -race -coverprofile=coverage.unit.out -covermode=atomic -v -coverpkg=./src/... ./src/...`
- `test_e2e` (depends on `test_clean`)
  - removes `coverage.e2e.out`
  - runs:
    - `ENV=test go test -count 1 -timeout 10s -race -coverprofile=coverage.e2e.out -covermode=atomic -v -coverpkg=./src/... ./tests/...`
- `test`
  - runs `test_unit` then `test_e2e`

## Environment mechanics to remember

- `ENV=test` triggers `src/config/env.go` logic:
  - changes cwd to repo root
  - loads `onepixel.test.env`
  - then loads `onepixel.local.env` as fallback defaults
- `onepixel.test.env` sets sqlite/duckdb dialects.
- `USE_FILE_DB` comes from fallback local env by default (`true` in `onepixel.local.env`).

## GeoIP/network caveat (important)

`db.GetGeoIPDB()` checks `./GeoLite2-City.mmdb` freshness (30 days). If stale/missing, it downloads from:
- `https://git.io/GeoLite2-City.mmdb`

In restricted/offline environments, tests touching GeoIP can fail due to network/DNS.

Mitigation before running tests:
- `curl -L https://git.io/GeoLite2-City.mmdb -o GeoLite2-City.mmdb`

## Recommended agent workflow

1. Run targeted tests for changed packages first.
2. Run `make test_unit`.
3. Run `make test_e2e` when API/integration behavior changed.
4. Use `make test` before final handoff when feasible.

## Troubleshooting notes

- If tests time out, remember Makefile uses `-timeout 10s`; for debugging, rerun failing package with a longer timeout manually.
- If local env has been modified to non-file DB settings, tests may accidentally try postgres/clickhouse; restore test-friendly env values.
- Coverage artifacts are written to:
  - `coverage.unit.out`
  - `coverage.e2e.out`
