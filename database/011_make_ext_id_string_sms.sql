-- +goose Up
ALTER TABLE sms
ALTER COLUMN external_id TYPE varchar(500);
-- +goose Down
ALTER TABLE sms
ALTER COLUMN external_id TYPE bigint;