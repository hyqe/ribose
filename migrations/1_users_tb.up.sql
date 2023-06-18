CREATE TABLE IF NOT EXISTS users (
    id BIGINT GENERATED ALWAYS AS IDENTITY,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    uuid UUID NOT NULL DEFAULT gen_random_uuid(),
    email TEXT NOT NULL
);
CREATE UNIQUE INDEX idx_users_uuid ON users (uuid);
CREATE UNIQUE INDEX idx_users_email ON users (email);