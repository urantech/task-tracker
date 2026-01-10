-- +goose Up
-- +goose StatementBegin
CREATE TABLE config (
    id BIGSERIAL PRIMARY KEY,
    job_name VARCHAR(30) NOT NULL UNIQUE
);

INSERT INTO config (job_name) VALUES ('daily-report');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS config;
-- +goose StatementEnd