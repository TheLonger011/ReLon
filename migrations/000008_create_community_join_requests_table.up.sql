CREATE TABLE community_join_requests (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    community_id UUID NOT NULL REFERENCES communities(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    created_at TIMESTAMP DEFAULT NOW(),
    UNIQUE (community_id,user_id)
);