CREATE TABLE settings (
  id TEXT PRIMARY KEY,
  value TEXT NOT NULL,
  previous_value TEXT NOT NULL DEFAULT '',
  setting_group TEXT NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_settings_group ON settings(setting_group);