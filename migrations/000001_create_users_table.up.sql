CREATE TABLE users(
    id UUID DEFAULT gen_random_uuid() PRIMARY KEY ,
    username VARCHAR(40) UNIQUE NOT NULL ,
    email VARCHAR(127) UNIQUE NOT NULL,
    password_hash VARCHAR(127) NOT NULL,
    karma INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW()
);
