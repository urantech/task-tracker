package common

import (
	"errors"
	"github.com/go-playground/validator/v10"
)

var (
	ErrInvalidRequest     = errors.New("invalid request")
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type ValidationError struct {
	Errors map[string]string
}

func (e *ValidationError) Error() string {
	return "validation failed"
}

func NewValidationError(ve validator.ValidationErrors) *ValidationError {
	return &ValidationError{
		Errors: mapValidationErrors(ve),
	}
}

func mapValidationErrors(ve validator.ValidationErrors) map[string]string {
	result := make(map[string]string)

	for _, fe := range ve {
		switch fe.Tag() {
		case "required":
			result[fe.Field()] = "is required"
		case "email":
			result[fe.Field()] = "must be a valid email"
		case "min":
			result[fe.Field()] = "is too short"
		case "task_status":
			result[fe.Field()] = "must be one of: TODO, IN_PROGRESS, DONE"
		default:
			result[fe.Field()] = "is invalid"
		}
	}

	return result
}
