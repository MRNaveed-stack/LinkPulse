-- Drop the global unique constraint on slug and add composite unique constraint
-- This allows each user to have links with different slugs, but prevents duplicate slugs within a user

-- Drop the global unique constraint on slug if present
ALTER TABLE links DROP CONSTRAINT IF EXISTS links_slug_key;

-- Drop composite constraint if it already exists before adding
ALTER TABLE links DROP CONSTRAINT IF EXISTS links_user_id_slug_unique;

-- Add a composite unique constraint on user_id and slug
ALTER TABLE links ADD CONSTRAINT links_user_id_slug_unique UNIQUE (user_id, slug);

