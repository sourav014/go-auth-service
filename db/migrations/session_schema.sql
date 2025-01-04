CREATE TABLE sessions (
    id VARCHAR(255) PRIMARY KEY NOT NULL,
    user_email VARCHAR(255) NOT NULL,
    refresh_token TEXT NOT NULL,
    is_revoked BOOLEAN DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL
);
