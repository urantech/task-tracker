package task

import (
	"backend-api/internal/common"
	"context"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/go-playground/validator/v10"
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

func (s *Service) CreateTask(ctx context.Context, req CreateRequest, userId int64) (Task, error) {
	if err := s.validator.Struct(req); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			return Task{}, common.NewValidationError(ve)
		}

		return Task{}, ErrInvalidRequest
	}

	task, err := s.storage.Create(ctx, req.Title, req.Description, userId)
	if err != nil {
		return Task{}, fmt.Errorf("create task: %w", err)
	}

	return task, nil
}
