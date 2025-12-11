-- +goose Up
CREATE TABLE IF NOT EXISTS sms(
    id     BIGSERIAL  PRIMARY KEY,
    external_id BIGINT NOT NULL,
    service_name VARCHAR(255) NOT NULL,
    service_id uuid NOT NULL,
    from_number    VARCHAR(255),
    number    VARCHAR(255),
    text      VARCHAR(2500),
    status    BIGINT,
    extended_status VARCHAR(255),
    cost       float8,
    date_created BIGINT,
    date_sent    BIGINT,
    raw          TEXT
);

-- +goose Down
DROP TABLE IF EXISTS sms