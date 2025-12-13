-- +goose Up
CREATE TABLE IF NOT EXISTS agent_info (
    id uuid PRIMARY KEY NOT NULL ,
    csr_id int NOT NULL REFERENCES certificate_requests(id),
    organization_id VARCHAR(255) NOT NULL,
    email  VARCHAR(255) NOT NULL,
    certificate_pem VARCHAR(65535) NOT NULL,
    is_active bool NOT NULL
);
-- +goose Down
DROP TABLE IF EXISTS agent_info;