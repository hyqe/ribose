CREATE TABLE IF NOT EXISTS users (
    id BIGINT GENERATED ALWAYS AS IDENTITY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    uuid UUID NOT NULL DEFAULT gen_random_uuid(),
    email TEXT NOT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS users_uuid_key ON users (uuid);
CREATE UNIQUE INDEX IF NOT EXISTS users_email_key ON users (email);