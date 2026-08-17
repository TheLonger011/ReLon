CREATE TABLE communities (
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY,
    name VARCHAR(63) NOT NULL UNIQUE,
    description VARCHAR(255),
    owner_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    is_private BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);