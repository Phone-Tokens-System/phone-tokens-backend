--+goose Up
CREATE TABLE IF NOT EXISTS user_profile(
    user_id uuid PRIMARY KEY NOT NULL REFERENCES users(id) ON DELETE CASCADE ,
    birth_date timestamp,
    age int,
    gender VARCHAR(255) NOT NULL ,
    country VARCHAR(255) NOT NULL ,
    region VARCHAR(255),
    city VARCHAR(255),
    education VARCHAR(255)
);

--+goose Down
DROP TABLE IF EXISTS user_profile;