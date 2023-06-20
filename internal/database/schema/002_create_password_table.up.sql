CREATE TABLE IF NOT EXISTS passwords (
    uuid UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    user_id BIGINT NOT NULL,
    salt TEXT NOT NULL,
    algorithm TEXT NOT NULL,
    hash TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS passwords_user_id_idx on passwords(user_id);