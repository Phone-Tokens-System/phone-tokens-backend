-- +goose Up
ALTER TABLE user_tokens
ADD COLUMN name VARCHAR(255);

-- +goose Down
ALTER TABLE user_tokens
DROP COLUMN name
