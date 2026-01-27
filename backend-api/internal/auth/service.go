package auth

import (
	"backend-api/internal/common"
	"backend-api/internal/user"
	"context"
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	myVal "backend-api/internal/validator"
)

type Service struct {
	userStorage *user.Storage
	signingKey  string
	validator   *validator.Validate
}

func NewService(userStorage *user.Storage, jwtSecret string) *Service {
	return &Service{
		userStorage: userStorage,
		signingKey:  jwtSecret,
		validator:   myVal.NewValidator(),
	}
}

func (s *Service) Login(ctx context.Context, req LoginRequest) (string, error) {
	if err := s.validator.Struct(req); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			return "", common.NewValidationError(ve)
		}

		return "", common.ErrInvalidRequest
	}

	u, err := s.userStorage.GetByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, common.ErrUserNotFound) {
			return "", common.ErrInvalidCredentials
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)); err != nil {
		return "", common.ErrInvalidCredentials
	}

	claims := jwt.MapClaims{
		"sub": u.Id,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(s.signingKey))
}
