package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

var (
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrUserNotFound      = errors.New("user not found")
)

type Storage struct {
	conn *sql.DB
}

func NewStorage(conn *sql.DB) *Storage {
	return &Storage{conn: conn}
}

func (s *Storage) Create(ctx context.Context, user User) (User, error) {
	const query = `
		INSERT INTO users (email, password, role_id)
		VALUES ($1, $2, $3)
		ON CONFLICT (email) DO NOTHING
		RETURNING id, email
	`

	var u User

	err := s.conn.
		QueryRowContext(ctx, query, user.Email, user.Password, user.RoleId).
		Scan(&u.Id, &u.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, ErrUserAlreadyExists
		}

		return User{}, fmt.Errorf("create user: %w", err)
	}

	return u, nil
}

func (s *Storage) GetByEmail(ctx context.Context, email string) (User, error) {
	const query = `
		SELECT id, email, password, role_id, created_at
		FROM users
		WHERE email = $1
	`

	var u User

	err := s.conn.
		QueryRowContext(ctx, query, email).
		Scan(&u.Id, &u.Email, &u.Password, &u.RoleId, &u.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, ErrUserNotFound
		}

		return User{}, fmt.Errorf("get user by email: %w", err)
	}

	return u, nil
}

func (s *Storage) GetRoleId(ctx context.Context, role UserRole) (int64, error) {
	const query = `
		SELECT id
		FROM roles
		WHERE name = $1
	`

	var id int64

	err := s.conn.QueryRowContext(ctx, query, role).Scan(&id)
	if err != nil {
		return 0, err
	}

	return id, nil
}
