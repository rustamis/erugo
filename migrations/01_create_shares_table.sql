CREATE TABLE IF NOT EXISTS shares (
  id INTEGER PRIMARY KEY AUTOINCREMENT,
  file_path TEXT,
  expiration_date TEXT,
  long_id TEXT,
  num_files INTEGER,
  total_size INTEGER,
  files TEXT
)