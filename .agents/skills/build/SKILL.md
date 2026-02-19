# Build Skill

Use this skill when building binaries, running the app, or preparing docs artifacts.

## Primary build command

- Preferred build command (skip Swagger regeneration):
  - `make build DOCS=false`
- This compiles `src/main.go` into `bin/onepixel` for current OS/arch.

## Make targets and behavior

Defined in `/home/runner/work/onepixel_backend/onepixel_backend/Makefile`:

- `docs`:
  - runs `swag init --pd -g server/server.go -d src --md src/docs -o src/docs`
  - regenerates files in `src/docs`
- `build`:
  - builds one binary with detected/current `GOOS`/`GOARCH`
- `build_all`:
  - produces linux/darwin/windows artifacts under `bin/`
- `run`:
  - depends on `build`, then executes `./bin/onepixel`
- `clean`:
  - `go clean` and removes `bin/*`

`DOCS=false` clears build dependency on `docs`; use it for faster/local CI-safe builds.

## Prerequisites

- Go toolchain compatible with `go.mod` (Go 1.22 target).
- Optional (only if generating docs):
  - `swag` CLI (`go install github.com/swaggo/swag/cmd/swag@latest`)
- Optional local dev hot-reload:
  - `air` (`go install github.com/air-verse/air@latest`)

## Build verification checklist

1. From repo root, run `make build DOCS=false`.
2. Confirm binary exists at `bin/onepixel`.
3. If API annotations changed, regenerate docs (`make docs`) and verify `src/docs/*` updates are intentional.

## Local runtime options

- `make run` for binary run flow.
- `air` for hot-reload using `.air.toml`.
- `docker-compose up` for containerized stack (app + postgres + clickhouse + metabase).

## Common pitfalls

- Running plain `make` executes first target (`docs`), which requires `swag`.
- Host-based routing needs correct env values (`ADMIN_SITE_HOST`, `MAIN_SITE_HOST`) to hit expected app path.
