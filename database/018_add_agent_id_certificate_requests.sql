-- +goose Up
ALTER TABLE certificate_requests ADD COLUMN agent_id uuid;
-- +goose Down
ALTER TABLE certificate_requests DROP COLUMN agent_id;