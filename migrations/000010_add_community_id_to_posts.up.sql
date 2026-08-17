ALTER TABLE posts ADD COLUMN community_id UUID REFERENCES communities(id) ON DELETE SET NULL;

CREATE INDEX idx_posts_community_id ON posts(community_id);