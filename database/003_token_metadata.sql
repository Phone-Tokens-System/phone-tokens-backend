-- +goose Up
ALTER TABLE user_tokens
    ADD COLUMN name TEXT NOT NULL DEFAULT 'default token',
    ADD COLUMN permissions TEXT[] NOT NULL DEFAULT ARRAY['sms', 'calls'],
    ADD COLUMN status TEXT NOT NULL DEFAULT 'active';

-- +goose Down
ALTER TABLE user_tokens
    DROP COLUMN IF EXISTS status,
    DROP COLUMN IF EXISTS permissions,
    DROP COLUMN IF EXISTS name;
