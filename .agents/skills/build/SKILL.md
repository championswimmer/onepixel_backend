---
name: build
description: Build and run onepixel binaries, generate Swagger docs when needed, and produce cross-platform release artifacts.
---

# Build Skill

Use this skill when you need to compile, run, or package the onepixel backend.

## Quick commands

- Fast local build (recommended default):
  - `make build DOCS=false`
- CI-like verbose build:
  - `make build DOCS=false ARGS="-v"`
- Build all release artifacts:
  - `make build_all DOCS=false ARGS="-v"`
- Run app via built binary:
  - `make run DOCS=false`
- Regenerate Swagger docs (only when API annotations changed):
  - `make docs`

## Current Makefile behavior

Defined in `Makefile` at repo root:

- `docs`
  - runs: `swag init --pd -g server/server.go -d src --md src/docs -o src/docs`
  - writes generated docs into `src/docs/`
- `build`
  - output: `bin/onepixel`
  - compiles `src/main.go`
  - default dependency: `docs` (unless `DOCS=false`)
- `build_all`
  - outputs:
    - `bin/onepixel-linux-amd64`
    - `bin/onepixel-darwin-amd64`
    - `bin/onepixel-darwin-arm64`
    - `bin/onepixel-windows-amd64.exe`
- `run`
  - depends on `build`, then executes `./bin/onepixel`
- `clean`
  - `go clean` and remove `bin/*`

## OS/arch resolution details

The Makefile normalizes platform values:

- OS: Darwin/macOS → `darwin`, Linux → `linux`, otherwise `windows`
- ARCH: `x86_64` → `amd64`, `i386` → `386`

You can still override by exporting `GOOS`/`GOARCH` before running `make`.

## Prerequisites

- Go 1.22 (matches `go.mod` and CI workflow)
- `swag` CLI only if docs generation is needed:
  - `go install github.com/swaggo/swag/cmd/swag@latest`
- Optional local dev hot reload:
  - `go install github.com/air-verse/air@latest`

## Recommended agent workflow

1. Build with `make build DOCS=false`.
2. Confirm binary exists at `bin/onepixel`.
3. Only run `make docs` if endpoint annotations or Swagger comments changed.
4. If docs changed intentionally, include updated files under `src/docs/`.

## Common pitfalls

- Running plain `make` executes the first target (`docs`) and fails if `swag` is missing.
- `run` also triggers docs unless `DOCS=false` is passed.
- If `GOOS`/`GOARCH` are set in the shell, builds may target an unexpected platform.
- Host-based routing at runtime depends on env vars (`ADMIN_SITE_HOST`, `MAIN_SITE_HOST`).
