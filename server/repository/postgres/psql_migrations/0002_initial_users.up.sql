-- +migrate Up
CREATE TABLE users (
    login VARCHAR(50) PRIMARY KEY,
    user_id VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL
);