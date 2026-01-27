-- +goose Up
-- +goose StatementBegin
CREATE TABLE config (
    id BIGSERIAL PRIMARY KEY,
    job_name VARCHAR(30) NOT NULL UNIQUE,
    duration INTEGER NOT NULL CHECK (duration > 0),
    execution_time TIME NOT NULL
);

INSERT INTO config (job_name, duration, execution_time)
VALUES ('daily-report', 1, '00:00');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS config;
-- +goose StatementEnd