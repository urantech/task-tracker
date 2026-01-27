-- +goose Up
-- +goose StatementBegin
CREATE TABLE outbox_emails (
    id BIGSERIAL PRIMARY KEY,
    payload JSONB NOT NULL,
    status TEXT NOT NULL CHECK (
        status IN ('NEW', 'SENT', 'FAILED')
    ),
    type TEXT NOT NULL CHECK (
        type IN ('WELCOME', 'REPORT')
    ),
    dedup_key TEXT UNIQUE NOT NULL,
    attempts INTEGER NOT NULL DEFAULT 0,
    next_attempt_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE INDEX idx_outbox_emails_status_next_attempt ON outbox_emails (status, next_attempt_at);
CREATE UNIQUE INDEX id_outbox_emails_dedup_key_type ON outbox_emails (dedup_key, type);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS outbox_emails;
-- +goose StatementEnd