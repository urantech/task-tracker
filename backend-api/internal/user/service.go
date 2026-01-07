package user

import (
	"backend-api/internal/common"
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidRequest = errors.New("invalid request")

type Service struct {
	storage   *Storage
	validator *validator.Validate
}

func NewService(storage *Storage) *Service {
	v := validator.New()

	v.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}

		return name
	})

	return &Service{
		storage:   storage,
		validator: v,
	}
}

type UserRole string

const (
	RoleUser  UserRole = "ROLE_USER"
	RoleAdmin UserRole = "ROLE_ADMIN"
)

func (s *Service) RegisterUser(ctx context.Context, req RegisterRequest) (RegisterResponse, error) {
	if err := s.validator.Struct(req); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			return RegisterResponse{}, common.NewValidationError(ve)
		}

		return RegisterResponse{}, ErrInvalidRequest
	}

	roleId, err := s.storage.GetRoleId(ctx, RoleUser)
	if err != nil {
		return RegisterResponse{}, fmt.Errorf("get default role: %w", err)
	}

	hash, err := bcrypt.GenerateFromPassword(
		[]byte(req.PasswordRaw),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return RegisterResponse{}, fmt.Errorf("hash password: %w", err)
	}

	var user User

	user.Email = req.Email
	user.Password = string(hash)
	user.RoleId = roleId

	createdUser, err := s.storage.Create(ctx, user)
	if err != nil {
		if errors.Is(err, ErrUserAlreadyExists) {
			return RegisterResponse{}, err
		}

		return RegisterResponse{}, fmt.Errorf("register user: %w", err)
	}

	return RegisterResponse{createdUser.Id, createdUser.Email}, nil
}
