ALTER TABLE IF EXISTS passwords DROP CONSTRAINT IF EXISTS passwords_user_id_idx;
DROP TABLE IF EXISTS passwords;