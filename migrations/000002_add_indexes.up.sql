CREATE UNIQUE INDEX idx_users_email
ON users(email);

CREATE INDEX idx_files_user_id
ON files(user_id);

CREATE INDEX idx_files_expires_at
ON files(expires_at);

CREATE INDEX idx_conversions_user_id
ON conversions(user_id);

CREATE INDEX idx_conversions_status
ON conversions(status);

CREATE INDEX idx_conversions_created_at
ON conversions(created_at DESC);

CREATE INDEX idx_jobs_status_created
ON jobs(status, created_at);
