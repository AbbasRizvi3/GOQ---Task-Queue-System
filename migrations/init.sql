\c goq_db;

CREATE TABLE IF NOT EXISTS tasks (
    id VARCHAR(8) PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    payload TEXT,
    state VARCHAR(50) NOT NULL DEFAULT 'pending',
    run_at TIMESTAMPTZ NOT NULL,
    next_run_at TIMESTAMPTZ NULL,
    lease_until TIMESTAMPTZ,
    max_retries INT DEFAULT 3,
    retries INT DEFAULT 0,
    error TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_tasks_state ON tasks (state);

CREATE INDEX IF NOT EXISTS idx_tasks_run_at ON tasks (run_at);

CREATE INDEX IF NOT EXISTS idx_tasks_next_run_at ON tasks (next_run_at);

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
);