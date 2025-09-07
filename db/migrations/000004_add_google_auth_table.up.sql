CREATE TABLE google_auth (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email TEXT UNIQUE NOT NULL,
    username VARCHAR NOT NULL,
    provider TEXT NOT NULL DEFAULT 'google', 
    provider_id VARCHAR NOT NULL, 
    role VARCHAR DEFAULT 'user',
    created_at TIMESTAMP DEFAULT now()
);
