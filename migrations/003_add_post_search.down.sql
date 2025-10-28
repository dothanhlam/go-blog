DROP INDEX IF EXISTS posts_title_tsv_idx;
DROP TRIGGER IF EXISTS posts_tsvector_update ON posts;
DROP FUNCTION IF EXISTS update_posts_tsvector();
ALTER TABLE posts DROP COLUMN IF EXISTS title_tsv;