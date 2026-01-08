package task

import (
	"backend-api/internal/common"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
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

func (s *Storage) UpdateTask(ctx context.Context, id, userId int64, req UpdateRequest) (Task, error) {
	if req.Title == nil && req.Description == nil && req.Status == nil {
		return Task{}, errors.New("nothing to update")
	}

	setParts := []string{}
	args := []any{}
	argIdx := 1

	if req.Title != nil {
		setParts = append(setParts, fmt.Sprintf("title = $%d", argIdx))
		args = append(args, *req.Title)
		argIdx++
	}

	if req.Description != nil {
		setParts = append(setParts, fmt.Sprintf("description = $%d", argIdx))
		args = append(args, *req.Description)
		argIdx++
	}

	if req.Status != nil {
		setParts = append(setParts, fmt.Sprintf("status = $%d", argIdx))
		args = append(args, *req.Status)
		argIdx++
	}

	query := fmt.Sprintf(`
		UPDATE tasks
		SET %s
		WHERE id = $%d AND user_id = $%d
		RETURNING id, title, description, status, user_id, created_at, updated_at
	`, strings.Join(setParts, ", "), argIdx, argIdx+1)

	args = append(args, id, userId)

	var task Task

	err := s.conn.QueryRowContext(ctx, query, args...).Scan(
		&task.Id,
		&task.Title,
		&task.Description,
		&task.Status,
		&task.UserId,
		&task.CreatedAt,
		&task.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Task{}, common.ErrTaskNotFound
		}

		return Task{}, fmt.Errorf("failed to scan task: %w", err)
	}

	return task, nil
}
