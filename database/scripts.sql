CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE user_projects (
    id SERIAL PRIMARY KEY,
    user_id UUID REFERENCES users(id),
    project_type VARCHAR(50) NOT NULL, -- 'personal' or 'org'
    username VARCHAR(100) NOT NULL,
    pat TEXT NOT NULL, -- Store encrypted PAT
    repo_names TEXT NOT NULL, -- Comma-separated repo names
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
