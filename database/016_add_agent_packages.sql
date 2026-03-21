-- +goose Up
CREATE TABLE IF NOT EXISTS agent_packages(
    id SERIAL PRIMARY KEY NOT NULL ,
    agent_id uuid NOT NULL REFERENCES agents(id),
    package_id uuid NOT NULL REFERENCES packages(id),
    status VARCHAR(50) NOT NULL,
    units_total int8 NOT NULL ,
    units_used int8 NOT NULL ,
    expires_at timestamp
);

-- +goose Down
DROP TABLE IF EXISTS agent_packages;