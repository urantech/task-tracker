package task

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

func (s *Storage) Create(ctx context.Context, title, description string, userId int64) (Task, error) {
	tx, err := s.conn.BeginTx(ctx, nil)
	if err != nil {
		return Task{}, fmt.Errorf("begin transaction: %w", err)
	}

	defer func() {
		if rbErr := tx.Rollback(); rbErr != nil {
			log.Printf("Error rollback transaction: %v", err)
		}
	}()

	const query = `
			INSERT INTO tasks (title, description, status, user_id)
			VALUES ($1, $2, $3, $4)
			RETURNING id, title, description, status, user_id, created_at, updated_at
	`

	var task Task

	err = tx.
		QueryRowContext(ctx, query, title, description, StatusTodo, userId).
		Scan(&task.Id, &task.Title, &task.Description, &task.Status, &task.UserId, &task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		return Task{}, fmt.Errorf("create task in db: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return Task{}, fmt.Errorf("commit transaction: %w", err)
	}

	return task, nil
}

func (s *Storage) GetAllUserTasks(ctx context.Context, userId int64) ([]Task, error) {
	const query = `
			SELECT id, title, description, status, user_id, created_at, updated_at
			FROM tasks
			WHERE user_id = $1
			ORDER BY updated_at DESC
	`

	rows, err := s.conn.QueryContext(ctx, query, userId)
	if err != nil {
		return nil, fmt.Errorf("error to query user tasks: %w", err)
	}

	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			log.Printf("Error close rows: %v", err)
		}
	}()

	tasks := make([]Task, 0)

	for rows.Next() {
		var task Task

		err := rows.Scan(
			&task.Id,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.UserId,
			&task.CreatedAt,
			&task.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error get user tasks: %w", err)
	}

	return tasks, nil
}
