package task

import (
	"backend-api/internal/common"
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/go-playground/validator/v10"
)

type Service struct {
	storage   *Storage
	validator *validator.Validate
}

func NewService(storage *Storage) *Service {
	v := common.NewValidator()

	err := v.RegisterValidation("task_status", func(fl validator.FieldLevel) bool {
		status := fl.Field().String()
		switch TaskStatus(status) {
		case StatusTodo, StatusInProgress, StatusDone:
			return true
		default:
			return false
		}
	})
	if err != nil {
		log.Printf("Failed to register validation: %v", err)
	}

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

		return Task{}, common.ErrInvalidRequest
	}

	task, err := s.storage.Create(ctx, req.Title, req.Description, userId)
	if err != nil {
		return Task{}, fmt.Errorf("create task: %w", err)
	}

	return task, nil
}

func (s *Service) GetUserTasks(ctx context.Context, userId int64) ([]Task, error) {
	tasks, err := s.storage.GetAllUserTasks(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("get all user tasks: %w", err)
	}

	return tasks, nil
}
