package outbox

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"
)

type Storage struct {
	conn *sql.DB
}

func NewStorage(conn *sql.DB) *Storage {
	return &Storage{conn: conn}
}

func (s *Storage) Save(ctx context.Context, outbox OutboxEmail) error {
	const query = `
		INSERT INTO outbox_emails (payload, status, dedup_key, type)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (dedup_key, type) DO NOTHING
	`

	_, err := s.conn.ExecContext(
		ctx,
		query,
		outbox.Payload,
		outbox.Status,
		outbox.DedupKey,
		outbox.Type,
	)
	if err != nil {
		return fmt.Errorf("db: %w", err)
	}

	return nil
}

func (s *Storage) FetchBatch(ctx context.Context, limit int) ([]OutboxEmail, error) {
	const query = `
		SELECT id, payload, attempts, type
		FROM outbox_emails
		WHERE status = 'NEW'
		  AND next_attempt_at <= now()
		ORDER BY created_at
		LIMIT $1
		FOR UPDATE SKIP LOCKED
	`

	rows, err := s.conn.QueryContext(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("query context: %w", err)
	}

	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Printf("Close rows error: %v", closeErr)
		}
	}()

	var result []OutboxEmail

	for rows.Next() {
		var email OutboxEmail

		if err := rows.Scan(&email.ID, &email.Payload, &email.Attempts, &email.Type); err != nil {
			return nil, fmt.Errorf("scan rows %w", err)
		}

		result = append(result, email)
	}

	return result, rows.Err()
}

func (s *Storage) MarkSent(ctx context.Context, id int64) error {
	const query = `
		UPDATE outbox_emails
		SET status = 'SENT',
		    updated_at = now()
		WHERE id = $1
	`

	_, err := s.conn.ExecContext(ctx, query, id)

	return err
}

func (s *Storage) MarkFailed(
	ctx context.Context,
	id int64,
	attempts int,
	nextAttempt time.Time,
) error {
	const query = `
		UPDATE outbox_emails
		SET status = 'FAILED',
		    attempts = $2,
		    next_attempt_at = $3,
		    updated_at = now()
		WHERE id = $1
	`

	_, err := s.conn.ExecContext(
		ctx,
		query,
		id,
		attempts,
		nextAttempt,
	)

	return err
}

func (s *Storage) MarkPermanentlyFailed(ctx context.Context, eventID int64) error {
	const query = `
		UPDATE outbox_emails
		SET status = 'FAILED',
		    updated_at = now()
		WHERE id = $1
	`

	_, err := s.conn.ExecContext(ctx, query, eventID)

	return err
}
