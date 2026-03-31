---
name: api-endpoints
description: Reference for all onepixel API endpoints, their HTTP methods, paths, auth requirements, and behavior.
---

# API Endpoints Skill

Use this skill when adding, modifying, or debugging API routes, middleware, or authentication logic.

## Authentication types

Four auth middleware are defined in `src/security/middlewares.go`:

- **MandatoryAdminApiKeyAuthMiddleware** – requires `X-API-Key` header matching `config.AdminApiKey`. Sets `admin=true` in Fiber locals.
- **OptionalAdminApiKeyAuthMiddleware** – same check, but proceeds without error if header is absent.
- **MandatoryJwtAuthMiddleware** – requires `Authorization: Bearer <token>`. Stores decoded `*models.User` in Fiber locals.
- **OptionalJwtAuthMiddleware** – same check, but proceeds without error if header is absent.

## Users endpoints (`/api/v1/users`)

Source: `src/routes/api/users.go`

| Method  | Path                | Auth                     | Description                                                        |
| ------- | ------------------- | ------------------------ | ------------------------------------------------------------------ |
| `POST`  | `/users`            | **X-API-Key** (required) | Register a new user. Returns JWT token. 409 if email exists.       |
| `POST`  | `/users/login`      | None                     | Login with email/password. Returns JWT token. 401 if invalid.      |
| `GET`   | `/users/{userid}`   | **JWT** (required)       | Get user info (currently returns placeholder).                     |
| `PATCH` | `/users/{userid}`   | **JWT** (required)       | Update user password. Returns updated user and new JWT token.      |

## URLs endpoints (`/api/v1/urls`)

Source: `src/routes/api/urls.go`

| Method  | Path                                      | Auth                                       | Description                                                                 |
| ------- | ----------------------------------------- | ------------------------------------------ | --------------------------------------------------------------------------- |
| `GET`   | `/urls`                                   | **JWT** (optional) + **X-API-Key** (optional) | Get all URLs. Admin can filter by `?userid=`. Returns 401 if neither auth.  |
| `GET`   | `/urls/groups`                            | **JWT** (required)                         | Get URL groups created by the authenticated user.                           |
| `POST`  | `/urls/groups`                            | **X-API-Key** (required)                   | Create a URL group. Requires `short_path` and `creator_id`.                 |
| `POST`  | `/urls/groups/{group}/shorten`            | **JWT** (required)                         | Create a random short URL inside a group the user owns.                     |
| `POST`  | `/urls/groups/{group}/shorten/{shortcode}`| **JWT** (required)                         | Create a specific short URL inside a group the user owns.                   |
| `GET`   | `/urls/groups/{group}/{shortcode}`        | None                                       | Get URL info (long URL + hit count) for a grouped short URL.               |
| `POST`  | `/urls`                                   | **JWT** (required)                         | Create a random short URL for the authenticated user.                       |
| `PUT`   | `/urls/{shortcode}`                       | **JWT** (required)                         | Create a short URL with a specific shortcode. 409 if exists, 403 if forbidden. |
| `GET`   | `/urls/{shortcode}`                       | None                                       | Get URL info (long URL + hit count). 404 if not found.                      |

## Stats endpoints (`/api/v1/stats`)

Source: `src/routes/api/stats.go`

| Method | Path                  | Auth                 | Description                                                         |
| ------ | --------------------- | -------------------- | ------------------------------------------------------------------- |
| `GET`  | `/stats`              | **JWT** (optional)   | Get redirect statistics for the authenticated user.                 |
| `GET`  | `/stats/{shortcode}`  | None                 | Get stats for a specific short URL (currently returns placeholder). |

## Redirect endpoints (Main app – `MAIN_SITE_HOST`)

Source: `src/routes/redirect/redirect.go`, `src/server/server.go`

| Method | Path                    | Auth | Description                                                         |
| ------ | ----------------------- | ---- | ------------------------------------------------------------------- |
| `GET`  | `/{group}/{shortcode}`  | None | Redirect to the long URL via group + shortcode. Logs event async.   |
| `GET`  | `/{shortcode}`          | None | Redirect to the long URL via shortcode. Logs event async.           |
| `GET`  | `/`                     | None | Redirect to admin host. Logs event async.                           |

## Static and documentation (Admin app)

Source: `src/server/server.go`

| Method | Path      | Auth | Description                                                  |
| ------ | --------- | ---- | ------------------------------------------------------------ |
| `GET`  | `/docs/*` | None | Swagger UI for OpenAPI documentation.                        |
| `GET`  | `/`       | None | Static files from `./public_html` (gzip, 1-hour cache).     |

## X-API-Key summary

Only these endpoints require or accept the `X-API-Key` header:

- `POST /api/v1/users` – **required** (admin-only user registration)
- `POST /api/v1/urls/groups` – **required** (admin-only group creation)
- `GET  /api/v1/urls` – **optional** (admin access to query any user's URLs)
