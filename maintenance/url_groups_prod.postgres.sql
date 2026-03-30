-- Production URL groups migration for Postgres.
--
-- Audit notes from 2026-03-30:
-- - no legacy global short_url index was present on public.urls
-- - default url group row (id = 0) already exists
-- - no duplicate url_groups.name rows exist
-- - no duplicate (url_group_id, short_url) rows exist
--
-- Run with:
--   psql "$DATABASE_URL" -v ON_ERROR_STOP=1 -f maintenance/url_groups_prod.postgres.sql

INSERT INTO public.url_groups (id, created_at, updated_at, deleted_at, name, creator_id)
VALUES (0, NOW(), NOW(), NULL, '0', 0)
ON CONFLICT (id) DO UPDATE
SET
    updated_at = EXCLUDED.updated_at,
    deleted_at = NULL,
    name = EXCLUDED.name,
    creator_id = EXCLUDED.creator_id;

DROP INDEX IF EXISTS public.idx_urls_short_url;
DROP INDEX IF EXISTS public.uni_urls_short_url;
DROP INDEX IF EXISTS public.short_url;
DROP INDEX IF EXISTS public."ShortURL";

CREATE UNIQUE INDEX IF NOT EXISTS idx_url_groups_short_path
    ON public.url_groups USING btree (name);

CREATE UNIQUE INDEX IF NOT EXISTS idx_urls_group_shortcode
    ON public.urls USING btree (url_group_id, short_url);
