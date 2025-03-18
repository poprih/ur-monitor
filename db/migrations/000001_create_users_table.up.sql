CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    line_user_id TEXT UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);
