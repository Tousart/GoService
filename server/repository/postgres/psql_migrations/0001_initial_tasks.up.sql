-- +migrate Up
CREATE TABLE tasks (
    task_id VARCHAR(255) PRIMARY KEY,
    status VARCHAR(50) NOT NULL,
    stdout VARCHAR(50),
    stderr VARCHAR(50)
);