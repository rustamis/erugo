CREATE TABLE IF NOT EXISTS share_access_logs (
  share_id INTEGER,
  user_email TEXT,
  user_ip TEXT,
  user_agent TEXT,
  access_date TEXT,
  PRIMARY KEY (share_id, user_email)
);