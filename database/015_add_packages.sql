--+goose Up
CREATE TABLE IF NOT EXISTS packages(
    id uuid PRIMARY KEY NOT NULL,
    name VARCHAR(255) NOT NULL UNIQUE ,
    service VARCHAR(50) NOT NULL ,
    units int8 NOT NULL,
    price float8 NOT NULL
);

--+goose Down
DROP TABLE IF EXISTS packages;