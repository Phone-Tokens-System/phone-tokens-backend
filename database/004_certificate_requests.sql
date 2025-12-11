-- +goose Up
CREATE TABLE IF NOT EXISTS certificate_requests (
    id SERIAL PRIMARY KEY ,
    email VARCHAR(255) NOT NULL,
    csr   VARCHAR(65535) NOT NULL ,
    status varchar(20) NOT NULL
);
-- +goose Down
DROP TABLE IF EXISTS certificate_requests