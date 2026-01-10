package job

import (
	"context"
	"database/sql"
	"fmt"
	"log"
)

type Storage struct {
	conn *sql.DB
}

func NewStorage(conn *sql.DB) *Storage {
	return &Storage{conn: conn}
}

func (s *Storage) GetAllJobConfigs(ctx context.Context) ([]Config, error) {
	const query = `
		SELECT job_name, duration, execution_time
		FROM config
		ORDER BY id
	`

	rows, err := s.conn.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("get all job configs: %w", err)
	}

	defer func() {
		if closeErr := rows.Close(); err != nil {
			log.Printf("failed to close rows: %v", closeErr)
		}
	}()

	configs := make([]Config, 0)

	for rows.Next() {
		var cfg Config

		err := rows.Scan(
			&cfg.JobName,
			&cfg.Duration,
			&cfg.ExecutionTime,
		)
		if err != nil {
			return nil, fmt.Errorf("scan job config: %w", err)
		}

		configs = append(configs, cfg)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return configs, nil
}
