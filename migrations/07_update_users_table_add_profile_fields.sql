ALTER TABLE users ADD COLUMN full_name TEXT;
ALTER TABLE users ADD COLUMN email TEXT;
ALTER TABLE users ADD COLUMN must_change_pw INTEGER DEFAULT 0;
ALTER TABLE users ADD COLUMN created_at DATETIME;
ALTER TABLE users ADD COLUMN updated_at DATETIME;
ALTER TABLE users ADD COLUMN active INTEGER DEFAULT 1;

UPDATE users SET full_name = username;
UPDATE users SET email = "no-email@example.com";
UPDATE users SET must_change_pw = 0;
UPDATE users SET created_at = CURRENT_TIMESTAMP;
UPDATE users SET updated_at = CURRENT_TIMESTAMP;
UPDATE users SET active = 1;
