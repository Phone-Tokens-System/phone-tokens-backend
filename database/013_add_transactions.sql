-- +goose Up
CREATE TABLE IF NOT EXISTS billing_transactions(
    id SERIAL PRIMARY KEY,
    agent_id UUID NOT NULL REFERENCES agents(id),
    amount float8 NOT NULL,
    type VARCHAR(20) NOT NULL ,
    service VARCHAR(20) NOT NULL ,
    created_at timestamp NOT NULL DEFAULT now(),
    stripe_session_id TEXT
);

-- +goose Down
DROP TABLE IF EXISTS billing_transactions;