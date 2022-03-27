DROP
DATABASE IF EXISTS dbUsers;
CREATE
DATABASE dbUsers;
USE
dbUsers;

DROP TABLE IF EXISTS users;
CREATE TABLE users
(
    id         serial PRIMARY KEY,
    email      VARCHAR(255) UNIQUE NOT NULL,
    firstName  VARCHAR(50)         NOT NULL,
    lastName   VARCHAR(50)         NOT NULL,
    password   VARCHAR(50)         NOT NULL,
    created_on TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

DROP TABLE IF EXISTS refresh_tokens;
CREATE TABLE refresh_tokens
(
    id         serial PRIMARY KEY,
    token      VARCHAR(255) UNIQUE NOT NULL,
    created_on TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO users (email, firstName, lastName, password)
VALUES ('example@example.com', 'First', 'Last', 'password');