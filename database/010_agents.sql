-- +goose Up
CREATE TABLE IF NOT EXISTS agents (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL UNIQUE REFERENCES users(id) ON DELETE CASCADE,
    service_name TEXT NOT NULL,
    email TEXT NOT NULL,
    certificate BYTEA NOT NULL DEFAULT ''::bytea,
    certificate_request BYTEA NOT NULL DEFAULT ''::bytea,
    balance NUMERIC NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE IF EXISTS agents;
