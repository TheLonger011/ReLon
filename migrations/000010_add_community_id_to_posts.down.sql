DROP INDEX IF EXISTS idx_posts_community_id;
ALTER TABLE posts DROP COLUMN community_id;