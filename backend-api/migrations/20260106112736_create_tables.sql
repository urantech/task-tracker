-- Active: 1767698927774@@127.0.0.1@5432@task-tracker-db
-- +goose Up
-- +goose StatementBegin
CREATE TYPE task_status AS ENUM ('TODO', 'IN_PROGRESS', 'DONE');

CREATE TYPE role_name AS ENUM ('ROLE_USER', 'ROLE_ADMIN');

CREATE TABLE roles (
    id BIGSERIAL PRIMARY KEY,
    name role_name UNIQUE NOT NULL
);

CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR NOT NULL,
    role_id BIGINT REFERENCES roles (id),
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE TABLE tasks (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    description VARCHAR,
    status task_status DEFAULT 'TODO',
    user_id BIGINT NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX idx_user_email ON users (email);

CREATE INDEX idx_tasks_user_status ON tasks (user_id, status);

CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_tasks_modtime
    BEFORE UPDATE ON tasks
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

INSERT INTO roles (name)
VALUES ('ROLE_USER'), ('ROLE_ADMIN');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS update_tasks_modtime ON tasks;

DROP FUNCTION IF EXISTS update_updated_at_column ();

DROP TABLE IF EXISTS tasks;

DROP TABLE IF EXISTS users;

DROP TABLE IF EXISTS roles;

DROP TYPE IF EXISTS role_name;

DROP TYPE IF EXISTS task_status;
-- +goose StatementEnd