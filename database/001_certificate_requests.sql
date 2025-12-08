CREATE TABLE IF NOT EXISTS certificate_requests(
    id int PRIMARY KEY NOT NULL CONSTRAINT SERIAL,
    email VARCHAR(255) NOT NULL,
    csr   VARCHAR(65535) NOT NULL ,
    status varchar(20) NOT NULL
)