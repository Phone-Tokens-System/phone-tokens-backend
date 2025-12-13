-- +goose Up
ALTER TABLE agent_info ADD CONSTRAINT unique_csr UNIQUE (csr_id);
-- +goose Down
ALTER TABLE agent_info DROP CONSTRAINT unique_csr