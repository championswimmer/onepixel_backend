# Plan: Frontend Dashboard for onepixel

## Problem Statement

Create a Vue.js + TypeScript + Vite frontend dashboard for the onepixel URL shortener.
The dashboard is **only for authenticated users** — there are no public/unauthenticated pages
(the existing landing page at `public_html/` already serves that role).

The dashboard will be hosted at **`app.onepixel.link`** via GitHub Pages (`gh-pages` branch),
and deployed automatically on pushes to `main` that change the `/dashboard` directory.

## Stack

| Layer         | Choice                                    |
|---------------|-------------------------------------------|
| Framework     | Vue 3 (Composition API + `<script setup>`) |
| Language      | TypeScript                                |
| Build tool    | Vite                                      |
| CSS framework | Bootstrap 5.3.2 + Bootstrap Icons         |
| Router        | Vue Router 4 (hash mode for GH Pages)    |
| HTTP client   | Fetch API (thin typed wrapper)            |
| Auth storage  | `localStorage` (JWT token)                |

### Design Notes

- **Spartan / minimal** — match the landing page aesthetic: Bootstrap defaults, no custom
  theme beyond what's in the landing page `style.css`, light/dark mode via
  `prefers-color-scheme` + `data-bs-theme`.
- Bootstrap loaded from CDN in `index.html` (consistent with landing page).
- No state management library — a simple reactive `auth` composable with `ref()` is enough.
- Hash-based routing (`createWebHashHistory`) avoids the need for server-side rewrites on
  GitHub Pages.

## API Surface (from Swagger spec)

Base URL (production): `https://onepixel.link/api/v1`

### Auth
| Method | Path               | Auth        | Purpose            |
|--------|--------------------|-------------|--------------------|
| POST   | `/users/login`     | None        | Login → JWT token  |
| GET    | `/users/{userid}`  | Bearer JWT  | Get user info      |
| PATCH  | `/users/{userid}`  | Bearer JWT  | Update password    |

### URLs
| Method | Path                                          | Auth       | Purpose                       |
|--------|-----------------------------------------------|------------|-------------------------------|
| GET    | `/urls`                                       | Bearer JWT | List user's URLs              |
| POST   | `/urls`                                       | Bearer JWT | Create random short URL       |
| PUT    | `/urls/{shortcode}`                           | Bearer JWT | Create custom short URL       |
| GET    | `/urls/{shortcode}`                           | Bearer JWT | Get URL info + hit count      |
| POST   | `/urls/groups`                                | API Key    | Create URL group (admin)      |
| POST   | `/urls/groups/{group}/shorten`                | Bearer JWT | Create grouped random URL     |
| POST   | `/urls/groups/{group}/shorten/{shortcode}`    | Bearer JWT | Create grouped custom URL     |

### Stats
| Method | Path                    | Auth       | Purpose            |
|--------|-------------------------|------------|--------------------|
| GET    | `/stats`                | Bearer JWT | Get user stats     |
| GET    | `/stats/{shortcode}`    | Bearer JWT | Per-URL stats (WIP)|

## Pages / Routes

| Route (hash)          | Component          | Purpose                                      |
|-----------------------|--------------------|----------------------------------------------|
| `#/login`             | `LoginView`        | Email + password form → stores JWT            |
| `#/`                  | `DashboardView`    | Overview: URL count, total hits, recent URLs  |
| `#/urls`              | `UrlsListView`     | Table of all user URLs with hit counts        |
| `#/urls/new`          | `CreateUrlView`    | Form: long URL + optional custom shortcode    |
| `#/urls/:shortcode`   | `UrlDetailView`    | Single URL detail + stats (when API ready)    |
| `#/account`           | `AccountView`      | Change password                               |

A **global navigation guard** redirects unauthenticated users to `#/login`.
After login, redirect to `#/`.

## Directory Structure

```
dashboard/
├── index.html              # Entry point (Bootstrap CDN, mount point)
├── package.json
├── tsconfig.json
├── tsconfig.node.json
├── vite.config.ts
├── env.d.ts
├── public/
│   └── favicon.ico
└── src/
    ├── main.ts             # App bootstrap
    ├── App.vue             # Root component (navbar + router-view)
    ├── router/
    │   └── index.ts        # Vue Router config + auth guard
    ├── composables/
    │   └── useAuth.ts      # Reactive auth state (token, user, login/logout)
    ├── api/
    │   ├── client.ts       # Typed fetch wrapper (base URL, auth header)
    │   ├── users.ts        # User API calls
    │   ├── urls.ts         # URL API calls
    │   └── stats.ts        # Stats API calls
    ├── types/
    │   └── index.ts        # TypeScript interfaces matching DTOs
    ├── views/
    │   ├── LoginView.vue
    │   ├── DashboardView.vue
    │   ├── UrlsListView.vue
    │   ├── CreateUrlView.vue
    │   ├── UrlDetailView.vue
    │   └── AccountView.vue
    └── components/
        ├── NavBar.vue      # Top nav with links + logout
        └── UrlTable.vue    # Reusable URL table component
```

## TypeScript Types (from API DTOs)

```typescript
// Request types
interface LoginRequest { email: string; password: string }
interface CreateUrlRequest { long_url: string }
interface CreateUrlGroupRequest { short_path: string; creator_id: number }
interface UpdateUserRequest { password: string }

// Response types
interface UserResponse { id: number; email: string; token: string }
interface UrlResponse { short_url: string; long_url: string; creator_id: number }
interface UrlInfoResponse { long_url: string; hit_count: number }
interface UrlGroupResponse { short_path: string; creator_id: number }
interface ErrorResponse { status: number; message: string }
```

## Auth Flow

1. User visits any `#/` route → guard checks `localStorage` for JWT token.
2. If no token → redirect to `#/login`.
3. Login form POSTs to `/api/v1/users/login` → receives `UserResponse` with `token`.
4. Token + user info stored in `localStorage`; reactive `useAuth()` composable updates.
5. All subsequent API calls include `Authorization: Bearer <token>` header.
6. Logout clears `localStorage` and redirects to `#/login`.
7. On 401 response from any API call, auto-logout and redirect to login.

## API Client Configuration

```typescript
const API_BASE = import.meta.env.VITE_API_BASE_URL || 'https://onepixel.link/api/v1'
```

- In dev: override via `.env.development` → `VITE_API_BASE_URL=http://localhost:3000/api/v1`
- In prod: defaults to `https://onepixel.link/api/v1`

## Dark/Light Theme Support

Match the landing page approach:
- Detect `prefers-color-scheme` via `window.matchMedia`
- Apply `data-bs-theme="dark"` or `data-bs-theme="light"` to `<html>`
- Listen for changes and update reactively

## GitHub Actions: Deploy Dashboard

New workflow file: `.github/workflows/deploy-dashboard.yaml`

```yaml
name: Deploy Dashboard

on:
  push:
    branches:
      - main
    paths:
      - 'dashboard/**'
  workflow_dispatch:   # allow manual trigger

jobs:
  deploy:
    runs-on: ubuntu-latest
    permissions:
      contents: write
    concurrency:
      group: deploy-dashboard-${{ github.ref }}
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-node@v4
        with:
          node-version: '20'
          cache: 'npm'
          cache-dependency-path: dashboard/package-lock.json

      - name: Install dependencies
        working-directory: dashboard
        run: npm ci

      - name: Build
        working-directory: dashboard
        run: npm run build

      - name: Deploy to GitHub Pages
        uses: peaceiris/actions-gh-pages@v4
        if: github.ref == 'refs/heads/main'
        with:
          github_token: ${{ secrets.GITHUB_TOKEN }}
          publish_dir: ./dashboard/dist
          cname: app.onepixel.link
          force_orphan: true
          user_name: 'github-actions[bot]'
          user_email: 'github-actions[bot]@users.noreply.github.com'
```

Key configuration:
- **`paths: ['dashboard/**']`** — only triggers when dashboard files change.
- **`cname: app.onepixel.link`** — automatically creates a `CNAME` file in gh-pages.
- **`force_orphan: true`** — keeps gh-pages branch clean (single commit).
- **`workflow_dispatch`** — allows manual re-deploy from Actions tab.

### First-time setup requirement
After the first deployment, the `gh-pages` branch must be selected as the Pages source
in the repository settings (Settings → Pages → Branch → `gh-pages`).

## Vite Config

```typescript
export default defineConfig({
  base: '/',              // root domain, not a subdirectory
  plugins: [vue()],
  server: {
    proxy: {
      '/api': {
        target: 'http://localhost:3000',
        changeOrigin: true,
      }
    }
  }
})
```

- `base: '/'` because `app.onepixel.link` is a root domain.
- Dev proxy forwards `/api` calls to the local Go backend.

## Implementation Todos

### Phase 1: Project Scaffolding
1. **`scaffold-vite-project`** — Initialize Vite + Vue + TS project in `/dashboard` using
   `npm create vite@latest`. Configure `vite.config.ts`, `tsconfig.json`, `package.json`.
2. **`setup-bootstrap`** — Add Bootstrap 5.3.2 + Bootstrap Icons CDN links in `index.html`.
   Add dark/light theme detection script.
3. **`setup-router`** — Install `vue-router`, create route definitions, add auth navigation guard.

### Phase 2: Core Infrastructure
4. **`create-types`** — Define TypeScript interfaces matching API DTOs in `src/types/index.ts`.
5. **`create-api-client`** — Build typed fetch wrapper with base URL config, auth header
   injection, and 401 auto-logout in `src/api/client.ts`.
6. **`create-api-modules`** — Implement `users.ts`, `urls.ts`, `stats.ts` API modules.
7. **`create-auth-composable`** — Build `useAuth()` composable: reactive token/user state,
   login/logout methods, localStorage persistence.

### Phase 3: Views and Components
8. **`create-navbar`** — NavBar component with links to Dashboard, URLs, Create URL, Account,
   and Logout button. Conditionally hidden on login page.
9. **`create-login-view`** — Login form with email/password, error display, redirect on success.
10. **`create-dashboard-view`** — Overview page: URL count, total hits summary, recent URLs table.
11. **`create-urls-list-view`** — Full table of user's URLs: shortcode, long URL, hit count,
    created date. Reusable `UrlTable` component.
12. **`create-url-detail-view`** — Single URL details page with stats (placeholder for future
    per-URL analytics).
13. **`create-url-form-view`** — Create URL form: long URL input, optional custom shortcode
    toggle, success feedback with copy-to-clipboard.
14. **`create-account-view`** — Change password form.

### Phase 4: Polish & Deploy
15. **`theme-support`** — Implement dark/light mode detection and toggle matching landing page.
16. **`app-shell`** — Wire up `App.vue` with NavBar + `<router-view>`, loading states,
    toast notifications for errors/success.
17. **`create-github-action`** — Add `.github/workflows/deploy-dashboard.yaml` with the
    peaceiris gh-pages action configuration.
18. **`verify-build`** — Run `npm run build` to ensure production build succeeds, check output
    in `dashboard/dist/`.

## Dependencies

| Phase | Todos | Depends On |
|-------|-------|------------|
| 1     | scaffold, bootstrap, router | — |
| 2     | types, api-client, api-modules, auth | Phase 1 |
| 3     | All views and components | Phase 2 |
| 4     | Theme, shell, action, verify | Phase 3 |

## Out of Scope

- User registration (admin-only via API key — not exposed in dashboard)
- URL group management (admin-only)
- Per-URL detailed analytics charts (API endpoint is WIP)
- PWA / offline support
- i18n / localization
