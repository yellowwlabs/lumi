CREATE TABLE IF NOT EXISTS users (
    id              UUID PRIMARY KEY,
    display_name    TEXT NOT NULL,
    username        TEXT NOT NULL UNIQUE,
    password_hash   TEXT NOT NULL,
    email           TEXT NOT NULL UNIQUE,
    email_verified  BOOLEAN NOT NULL DEFAULT FALSE,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT now()
);
