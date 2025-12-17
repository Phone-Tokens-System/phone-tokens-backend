-- +goose Up
ALTER TABLE user_tokens
ADD COLUMN agent_id uuid;

-- +goose Down
ALTER TABLE user_tokens
DROP agent_id;