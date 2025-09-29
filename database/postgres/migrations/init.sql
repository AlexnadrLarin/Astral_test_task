CREATE TABLE IF NOT EXISTS documents (
    id           UUID PRIMARY KEY,
    name         TEXT NOT NULL UNIQUE,
    mime         TEXT NOT NULL,
    file         BOOLEAN NOT NULL DEFAULT false,
    public       BOOLEAN NOT NULL DEFAULT false,
    owner_login  TEXT,               
    grant_list   TEXT[] DEFAULT '{}',
    created_at   TIMESTAMP NOT NULL DEFAULT NOW(),
    json_data    JSONB,
    file_path    TEXT
);

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    login TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS sessions (
    token      TEXT PRIMARY KEY,
    user_id    INTEGER NOT NULL,
    login      TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);