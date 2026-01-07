package auth

import (
	"backend-api/internal/common"
	"backend-api/internal/user"
	"context"
	"errors"
	"reflect"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrInvalidRequest     = errors.New("invalid request")
)

type Service struct {
	userStorage *user.Storage
	signingKey  string
	validator   *validator.Validate
}

func NewService(userStorage *user.Storage, jwtSecret string) *Service {
	v := validator.New()

	v.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]
		if name == "-" {
			return ""
		}

		return name
	})

	return &Service{
		userStorage: userStorage,
		signingKey:  jwtSecret,
		validator:   v,
	}
}

func (s *Service) Login(ctx context.Context, req LoginRequest) (string, error) {
	if err := s.validator.Struct(req); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			return "", common.NewValidationError(ve)
		}

		return "", ErrInvalidRequest
	}

	u, err := s.userStorage.GetByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, user.ErrUserNotFound) {
			return "", ErrInvalidCredentials
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password)); err != nil {
		return "", ErrInvalidCredentials
	}

	claims := jwt.MapClaims{
		"sub": u.Id,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
		"iat": time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(s.signingKey))
}
