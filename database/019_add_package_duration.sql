-- +goose Up
ALTER TABLE packages ADD COLUMN IF NOT EXISTS duration_days INT NOT NULL DEFAULT 30;
ALTER TABLE packages ADD COLUMN IF NOT EXISTS description TEXT;
ALTER TABLE agent_packages ADD COLUMN IF NOT EXISTS service_type VARCHAR(50) NOT NULL DEFAULT 'SMS';

-- +goose Down
ALTER TABLE packages DROP COLUMN IF EXISTS duration_days;
ALTER TABLE packages DROP COLUMN IF EXISTS description;
ALTER TABLE agent_packages DROP COLUMN IF EXISTS service_type;
