CREATE TABLE sessions (
    id UUID PRIMARY KEY,
    username VARCHAR UNIQUE NOT NULL,
    refresh_token VARCHAR NOT NULL,
    user_agent VARCHAR NOT NULL,
    client_ip VARCHAR NOT NULL,
    is_blocked BOOLEAN NOT NULL DEFAULT FALSE,
    expires_at TIMESTAMPTZ NOT NULL
);

ALTER TABLE sessions
ADD CONSTRAINT fk_sessions_username
FOREIGN KEY (username)
REFERENCES users (username)
ON DELETE CASCADE;