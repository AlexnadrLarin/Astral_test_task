CREATE TABLE IF NOT EXISTS documents (
    id          UUID PRIMARY KEY,
    name        TEXT NOT NULL UNIQUE,
    mime        TEXT NOT NULL,
    file        BOOLEAN NOT NULL DEFAULT false,
    public      BOOLEAN NOT NULL DEFAULT false,
    token       TEXT NOT NULL,              
    grant_list  TEXT[] DEFAULT '{}',        
    created_at  TIMESTAMP NOT NULL DEFAULT NOW(),
    json_data   JSONB,
    file_path   TEXT UNIQUE
);