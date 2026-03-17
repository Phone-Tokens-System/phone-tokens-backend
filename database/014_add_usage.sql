-- +goose Up
CREATE TABLE IF NOT EXISTS usage(
    id SERIAL PRIMARY KEY,
    agent_id UUID REFERENCES agents(id) NOT NULL,
    phone_number VARCHAR(50),
    service VARCHAR(20),
    units int,
    cost float8,
    created_at timestamp NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE IF EXISTS usage;