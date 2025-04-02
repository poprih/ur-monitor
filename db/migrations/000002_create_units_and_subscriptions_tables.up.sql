CREATE TABLE units (
    id SERIAL PRIMARY KEY,
    unit_name TEXT UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE subscriptions (
    id SERIAL PRIMARY KEY,
    line_user_id TEXT NOT NULL REFERENCES users(line_user_id) ON DELETE CASCADE,
    unit_id INTEGER NOT NULL REFERENCES units(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT NOW(),
    deleted_at TIMESTAMP DEFAULT NULL,
    UNIQUE (line_user_id, unit_id)
);

ALTER TABLE users ADD COLUMN reply_token TEXT;
ALTER TABLE users ADD COLUMN active BOOLEAN DEFAULT TRUE;
