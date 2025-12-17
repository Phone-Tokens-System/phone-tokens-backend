-- +goose Up
ALTER TABLE agent_info
ADD service_name VARCHAR(255);
ALTER TABLE certificate_requests
ADD service_name VARCHAR(255);

-- +goose Down
ALTER TABLE agent_info
DROP service_name;

ALTER TABLE certificate_requests
DROP service_name;
