package task

import (
	"context"
	"database/sql"
	"fmt"
)

type Storage struct {
	conn *sql.DB
}

func NewStorage(conn *sql.DB) *Storage {
	return &Storage{conn: conn}
}

func (s *Storage) Create(ctx context.Context, title, description string, userId int64) (Task, error) {
	const query = `
			INSERT INTO tasks (title, description, status, user_id)
			VALUES ($1, $2, $3, $4)
			RETURNING id, title, description, status, user_id, created_at, updated_at
	`
	var task Task

	err := s.conn.
		QueryRowContext(ctx, query, title, description, StatusTodo, userId).
		Scan(&task.Id, &task.Title, &task.Description, &task.Status, &task.UserId, &task.CreatedAt, &task.UpdatedAt)
	if err != nil {
		return Task{}, fmt.Errorf("create task in db: %w", err)
	}

	return task, nil
}
