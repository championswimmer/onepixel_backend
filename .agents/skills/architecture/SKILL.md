---
name: architecture
description: Understand onepixel runtime architecture, host-based app dispatch, and redirect analytics flow before changing behavior across routes/controllers/db.
---

# Architecture Skill

Use this skill when implementing features that cross app startup, routing, controller behavior, or data boundaries.

## Runtime model

- Process entrypoint: `src/main.go`.
- Startup initializes three singleton data clients:
  - app DB (`db.GetAppDB`)
  - events DB (`db.GetEventsDB`)
  - GeoIP DB (`db.GetGeoIPDB`)
- HTTP entry app dispatches by `Host` header:
  - `ADMIN_SITE_HOST` -> `server.CreateAdminApp`
  - `MAIN_SITE_HOST` -> `server.CreateMainApp`

## Data boundaries

- App DB stores users, URL groups, URLs.
- Events DB stores redirect analytics events.
- DB singletons and migrations are in `src/db/init.go`.

## Redirect and analytics flow

- Redirect handlers are in `src/routes/redirect/redirect.go`.
- Redirect resolution logs events asynchronously using `EventsController.LogRedirectAsync`.
- Event enrichment includes GeoIP and PostHog emission (`src/utils/posthog/posthog_client.go`).

## API layering pattern

- Preserve this stack for new/changed endpoints:
  - route handler -> parser (`src/server/parsers`) -> validator (`src/server/validators`) -> controller -> DTO response

